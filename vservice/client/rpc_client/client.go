package rpc_client

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	opentrace_go_grpc "github.com/Snowlights/gogrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"time"
	"vtool/parse"
	"vtool/vconfig"
	"vtool/vlog"
	clientCommon "vtool/vservice/client/common"
	"vtool/vservice/client/pool"
	"vtool/vservice/common"
	"vtool/vservice/server"
	"vtool/vtrace"
)

type RpcClient struct {
	client common.Client

	servType common.ServiceType

	thriftClient func(thrift.TTransport, thrift.TProtocolFactory) interface{}
	grpcClient   func(conn *grpc.ClientConn) interface{}

	clientPool *pool.ClientPool

	center vconfig.Center

	mu           sync.RWMutex
	clientConfig *vconfig.ClientConfig
}

func NewRpcClient(client common.Client, thriftServCli func(thrift.TTransport, thrift.TProtocolFactory) interface{}, grpcServCli func(conn *grpc.ClientConn) interface{}) common.RpcClient {

	rpcCli := &RpcClient{
		client:       client,
		thriftClient: thriftServCli,
		grpcClient:   grpcServCli,
	}

	newConnFunc := rpcCli.newThriftConn
	rpcCli.servType = common.Thrift
	if grpcServCli != nil {
		newConnFunc = rpcCli.newGrpcConn
		rpcCli.servType = common.Grpc
	}

	err := rpcCli.initCenter()
	if err != nil {
		vlog.ErrorF(context.Background(), "initCenter error: %v", err)
	}

	rpcCli.center.AddListener(&common.ClientListener{Change: rpcCli.reload})
	rpcCli.reload()
	cfg := rpcCli.getConfig()
	cfg = vconfig.DefaultClientConfig
	rpcCli.clientPool = pool.NewClientPool(&pool.ClientPoolConfig{
		ServiceName:    client.ServName(),
		Idle:           cfg.Idle,
		Active:         cfg.MaxActive,
		IdleTimeout:    time.Duration(cfg.IdleTimeout) * time.Millisecond,
		Wait:           cfg.Wait,
		WaitTimeOut:    time.Duration(cfg.WaitTimeout) * time.Millisecond,
		StatTime:       time.Duration(cfg.StatTime) * time.Millisecond,
		GetConnTimeout: time.Duration(cfg.GetConnTimeout) * time.Millisecond,
	}, newConnFunc)
	rpcCli.client.AddPoolHandler(rpcCli.deleteAddrHandler)
	return rpcCli
}

func (c *RpcClient) Rpc(ctx context.Context, args *common.ClientCallerArgs, fnRpc func(context.Context, interface{}) error) error {

	if len(args.HashKey) == 0 {
		args.HashKey = clientCommon.NewHashKey()
	}

	if args.TimeOut == 0 {
		args.TimeOut = time.Duration(c.getConfig().CallTimeout) * time.Millisecond
	}
	if args.TimeOut == 0 {
		args.TimeOut = vconfig.DefaultCallTimeout
	}

	serv, ok := c.client.GetServAddr(args.Lane, c.servType, args.HashKey)
	if !ok {
		return fmt.Errorf("%s caller args is %+v", common.NotFoundServInfo, args)
	}
	if serv.Type != common.Rpc {
		return fmt.Errorf("%s serv info is %+v, caller args is %+v", common.NotFoundServEngine, serv, args)
	}

	tCtx, cancel := context.WithTimeout(ctx, args.TimeOut)
	defer cancel()

	return c.do(tCtx, serv, fnRpc)
}

func (c *RpcClient) do(ctx context.Context, serv *common.ServiceInfo, fnRpc func(context.Context, interface{}) error) error {
	var err error
	retry := c.getConfig().RetryTime
	if retry == 0 {
		retry = vconfig.CallRetryTimes
	}
	for ; retry >= 0; retry-- {
		err = c.rpc(ctx, serv, fnRpc)
		if err == nil {
			return nil
		}
	}
	return err
}

func (c *RpcClient) rpc(ctx context.Context, serv *common.ServiceInfo, fnRpc func(context.Context, interface{}) error) error {
	conn, err := c.clientPool.Get(ctx, serv)
	if err != nil {
		return err
	}
	defer c.clientPool.Put(ctx, serv, conn)

	vtrace.SpanFromContent(ctx)
	return fnRpc(ctx, conn.GetConn())
}

func (c *RpcClient) updateConfig(cfg *vconfig.ClientConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.clientConfig = cfg
}

func (c *RpcClient) getConfig() *vconfig.ClientConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cfg := c.clientConfig
	return cfg
}

func (c *RpcClient) reload() {

	cfg := new(vconfig.ClientConfig)
	err := c.center.UnmarshalWithNameSpace(vconfig.Client, parse.PropertiesTagName, cfg)
	if err != nil {
		return
	}
	c.updateConfig(cfg)
	c.resetPoolConfig(cfg)
}

// todo: might have some problem, like auth, use secret key to fix it
func (c *RpcClient) initCenter() error {
	cfg, err := c.parseConfigEnv()
	if err != nil {
		return err
	}

	center, err := vconfig.NewCenter(cfg)
	if err != nil {
		return err
	}

	c.center = center
	return nil
}

func (c *RpcClient) parseConfigEnv() (*vconfig.CenterConfig, error) {
	centerConfig, err := vconfig.ParseConfigEnv()
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(centerConfig.Port, 10, 64)
	if err != nil {
		return nil, err
	}

	return &vconfig.CenterConfig{
		AppID:            c.client.ServGroup() + common.Slash + c.client.ServName(),
		Cluster:          centerConfig.Cluster,
		Namespace:        []string{vconfig.Client},
		IP:               centerConfig.IP,
		Port:             int(port),
		IsBackupConfig:   false,
		BackupConfigPath: "",
		MustStart:        centerConfig.MustStart,
	}, nil
}

func (c *RpcClient) resetPoolConfig(cfg *vconfig.ClientConfig) {
	if c.clientPool == nil {
		return
	}

	c.clientPool.ResetConfig(cfg)
	c.clientPool.ResetConnConfig(cfg)
}

func (c *RpcClient) deleteAddrHandler(addr []string) {
	for _, addr := range addr {
		c.clientPool.Delete(context.Background(), addr)
	}
}

func (c *RpcClient) newThriftConn(addr string) (common.RpcConn, error) {
	ctx := context.Background()

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(addr)
	if err != nil {
		vlog.ErrorF(ctx, "open thriftSocket addr:%s, err:%v", addr, err)
		return nil, err
	}
	useTransport := transportFactory.GetTransport(transport)

	if err := useTransport.Open(); err != nil {
		vlog.ErrorF(ctx, "open addr:%s err:%v", addr, err)
		return nil, err
	}

	return &ThriftClientConn{
		thriftSocket:  transport,
		transport:     useTransport,
		serviceClient: c.thriftClient(useTransport, protocolFactory),
	}, nil
}

func (c *RpcClient) newGrpcConn(addr string) (common.RpcConn, error) {
	ctx := context.Background()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(opentrace_go_grpc.OpenTracingClientInterceptor(opentracing.GlobalTracer(),
			opentrace_go_grpc.SpanDecorator(server.SpanDecorator))),
		grpc.WithStreamInterceptor(opentrace_go_grpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer(),
			opentrace_go_grpc.SpanDecorator(server.SpanDecorator))),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		vlog.ErrorF(ctx, "dial grpc addr: %s failed, err: %v", addr, err)
		return nil, err
	}

	return &GrpcClientConn{
		serviceClient: c.grpcClient(conn),
		conn:          conn,
	}, nil
}
