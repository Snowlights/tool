package rpc_client

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	opentrace_go_grpc "github.com/Snowlights/gogrpc"
	"github.com/Snowlights/tool/breaker"
	"github.com/Snowlights/tool/parse"
	"github.com/Snowlights/tool/vconfig"
	"github.com/Snowlights/tool/vlog"
	clientCommon "github.com/Snowlights/tool/vservice/client/common"
	"github.com/Snowlights/tool/vservice/client/pool"
	"github.com/Snowlights/tool/vservice/common"
	"github.com/Snowlights/tool/vservice/server"
	"github.com/Snowlights/tool/vtrace"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const clientFuncIndex = 5

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
	rpcCli.clientPool = pool.NewClientPool(rpcCli.loadClientPoolConfig(client.ServName()), newConnFunc)
	rpcCli.client.AddPoolHandler(rpcCli.deleteAddrHandler)

	err = breaker.InitBreakerManager(rpcCli.client.ServGroup(), rpcCli.client.ServName())
	if err != nil {
		err = nil
	}

	return rpcCli
}

func (c *RpcClient) loadClientPoolConfig(servName string) *pool.ClientPoolConfig {
	clientConfig, apolloConfig := vconfig.DefaultClientConfig, c.getConfig()
	clientPoolConfig := &pool.ClientPoolConfig{
		ServiceName:    servName,
		Idle:           clientConfig.Idle,
		Active:         clientConfig.MaxActive,
		IdleTimeout:    time.Duration(clientConfig.IdleTimeout) * time.Millisecond,
		Wait:           clientConfig.Wait,
		WaitTimeOut:    time.Duration(clientConfig.WaitTimeout) * time.Millisecond,
		StatTime:       time.Duration(clientConfig.StatTime) * time.Millisecond,
		GetConnTimeout: time.Duration(clientConfig.GetConnTimeout) * time.Millisecond,
	}

	clientPoolConfig.Wait = apolloConfig.Wait
	if apolloConfig.StatTime != 0 {
		clientPoolConfig.StatTime = time.Duration(apolloConfig.StatTime) * time.Millisecond
	}
	if apolloConfig.IdleTimeout != 0 {
		clientPoolConfig.IdleTimeout = time.Duration(apolloConfig.IdleTimeout) * time.Millisecond
	}
	if apolloConfig.WaitTimeout != 0 {
		clientPoolConfig.WaitTimeOut = time.Duration(apolloConfig.WaitTimeout) * time.Millisecond
	}
	if apolloConfig.GetConnTimeout != 0 {
		clientPoolConfig.GetConnTimeout = time.Duration(apolloConfig.GetConnTimeout) * time.Millisecond
	}
	if apolloConfig.MaxActive != 0 {
		clientPoolConfig.Active = apolloConfig.MaxActive
	}
	if apolloConfig.Idle != 0 {
		clientPoolConfig.Idle = apolloConfig.Idle
	}

	return clientPoolConfig
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

	ctx = c.injectServ(ctx)
	vtrace.SpanFromContent(ctx)

	if !breaker.Entry(c.client.ServName(), c.GetFuncName()) {
		return breaker.ErrTriggerBreaker(c.client.ServName(), c.GetFuncName())
	}

	err = fnRpc(ctx, conn.GetConn())
	breaker.StatBreaker(c.client.ServName(), c.GetFuncName(), err)

	return err
}

func (c *RpcClient) GetFuncName() string {
	funcName := ""
	pc, _, _, ok := runtime.Caller(clientFuncIndex)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
		if index := strings.LastIndex(funcName, "."); index != -1 {
			if len(funcName) > index+1 {
				funcName = funcName[index+1:]
			}
		}
	}
	return funcName
}

func (c *RpcClient) injectServ(ctx context.Context) context.Context {
	servBase := server.GetServBase()
	if servBase != nil {
		servInfo := servBase.ServInfo()
		if servInfo != nil {
			injectServInfoMap := make(map[string]string)
			injectServInfoMap[vtrace.Lane] = servInfo.Lane
			serv, ok := servInfo.ServList[common.Grpc]
			if ok {
				injectServInfoMap[vtrace.ServType] = string(common.Grpc)
				injectServInfoMap[vtrace.ServIP] = serv.Addr
				injectServInfoMap[vtrace.ServIP] = string(serv.Type)
			}
			serv, ok = servInfo.ServList[common.Thrift]
			if ok {
				injectServInfoMap[vtrace.ServType] = string(common.Thrift)
				injectServInfoMap[vtrace.ServIP] = serv.Addr
				injectServInfoMap[vtrace.EngineType] = string(serv.Type)
			}
			ctx = context.WithValue(ctx, clientCommon.InjectServKey, injectServInfoMap)
		}
	}
	return ctx
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
