package grpc

import (
	"context"
	"fmt"
	opentrace_go_grpc "github.com/Snowlights/gogrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"vtool/parse"
	"vtool/vconfig"
	"vtool/vlog"
	clientCommon "vtool/vservice/client/common"
	"vtool/vservice/client/pool"
	"vtool/vservice/common"
	"vtool/vservice/server"
)

type GrpcClient struct {
	client common.Client

	serviceClient func(conn *grpc.ClientConn) interface{}

	clientPool *pool.ClientPool

	center vconfig.Center

	mu   sync.RWMutex
	conf *vconfig.ClientConfig
}

func NewGrpcClient(client common.Client, servCli func(conn *grpc.ClientConn) interface{}) common.RpcClient {

	gc := &GrpcClient{
		client:        client,
		serviceClient: servCli,
		conf: &vconfig.ClientConfig{
			Idle:           pool.DefaultIdle,
			IdleTimeout:    pool.DefaultIdleTimeout,
			MaxActive:      pool.DefaultMaxActive,
			StatTime:       pool.DefaultStatTime,
			Wait:           true,
			WaitTimeout:    pool.DefaultWaitTimeout,
			GetConnTimeout: pool.DefaultGetConnTimeout,
		},
	}

	err := gc.initCenter()
	if err != nil {
		vlog.ErrorF(context.Background(), "initCenter error: %v", err)
	}

	gc.center.AddListener(&common.ClientListener{Change: gc.reload})
	gc.reload()
	cfg := gc.getConfig()

	gc.clientPool = pool.NewClientPool(&pool.ClientPoolConfig{
		ServiceName:    client.ServName(),
		Idle:           cfg.Idle,
		Active:         cfg.MaxActive,
		IdleTimeout:    cfg.IdleTimeout,
		Wait:           cfg.Wait,
		WaitTimeOut:    cfg.WaitTimeout,
		StatTime:       cfg.StatTime,
		GetConnTimeout: cfg.GetConnTimeout,
	}, gc.newConn)
	gc.client.AddPoolHandler(gc.deleteAddrHandler)
	return gc
}

func (g *GrpcClient) Rpc(args *common.ClientCallerArgs, fnRpc func(interface{}) error) error {
	if len(args.HashKey) == 0 {
		args.HashKey = clientCommon.NewHashKey()
	}

	serv, ok := g.client.GetServAddr(args.Lane, common.Grpc, args.HashKey)
	if !ok {
		return fmt.Errorf("%s caller args is %+v", common.NotFoundServInfo, args)
	}
	if serv.Type != common.Rpc {
		return fmt.Errorf("%s serv info is %+v, caller args is %+v", common.NotFoundServEngine, serv, args)
	}

	return g.do(context.TODO(), serv, fnRpc)
}

func (g *GrpcClient) updateConfig(cfg *vconfig.ClientConfig) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.conf = cfg
}

func (g *GrpcClient) getConfig() *vconfig.ClientConfig {
	g.mu.RLock()
	defer g.mu.RUnlock()

	cfg := g.conf
	return cfg
}

func (g *GrpcClient) reload() {

	cfg := new(vconfig.ClientConfig)
	err := g.center.UnmarshalWithNameSpace(vconfig.Client, parse.PropertiesTagName, cfg)
	if err != nil {
		return
	}
	go g.updateConfig(cfg)
	go g.resetPoolConfig(cfg)
}

func (g *GrpcClient) resetPoolConfig(cfg *vconfig.ClientConfig) {
	if g.clientPool == nil {
		return
	}
	g.clientPool.ResetConnConfig(cfg)
}

// todo: might have some problem, like auth, use secret key to fix it
func (g *GrpcClient) initCenter() error {
	cfg, err := g.parseConfigEnv()
	if err != nil {
		return err
	}

	center, err := vconfig.NewCenter(cfg)
	if err != nil {
		return err
	}

	g.center = center
	return nil
}

func (g *GrpcClient) parseConfigEnv() (*vconfig.CenterConfig, error) {
	centerConfig, err := vconfig.ParseConfigEnv()
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(centerConfig.Port, 10, 64)
	if err != nil {
		return nil, err
	}

	return &vconfig.CenterConfig{
		AppID:            g.client.ServGroup() + common.Slash + g.client.ServName(),
		Cluster:          centerConfig.Cluster,
		Namespace:        []string{vconfig.Client},
		IP:               centerConfig.IP,
		Port:             int(port),
		IsBackupConfig:   false,
		BackupConfigPath: "",
		MustStart:        centerConfig.MustStart,
	}, nil
}

func (g *GrpcClient) do(ctx context.Context, serv *common.ServiceInfo, fnRpc func(interface{}) error) error {
	var err error
	retry := 3
	for ; retry >= 0; retry-- {
		err = g.rpc(ctx, serv, fnRpc)
		if err == nil {
			return nil
		}
	}
	return err
}

func (g *GrpcClient) rpc(ctx context.Context, serv *common.ServiceInfo, fnRpc func(interface{}) error) error {
	conn, err := g.clientPool.Get(ctx, serv)
	if err != nil {
		return err
	}
	defer g.clientPool.Put(ctx, serv, conn)

	return fnRpc(conn.GetConn())
}

func (g *GrpcClient) deleteAddrHandler(addr []string) {
	for _, addr := range addr {
		g.clientPool.Delete(context.Background(), addr)
	}
}

func (g *GrpcClient) newConn(addr string) (common.RpcConn, error) {
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

	return &ClientConn{
		serviceClient: g.serviceClient(conn),
		conn:          conn,
	}, nil
}
