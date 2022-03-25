package grpc

import (
	"context"
	"fmt"
	opentrace_go_grpc "github.com/Snowlights/gogrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
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
}

func NewGrpcClient(client common.Client, servCli func(conn *grpc.ClientConn) interface{}) common.RpcClient {

	gc := &GrpcClient{
		client:        client,
		serviceClient: servCli,
	}
	gc.clientPool = pool.NewClientPool(&pool.ClientPoolConfig{
		ServiceName: client.ServName(),
		Idle:        pool.DefaultIdle,
		Active:      pool.DefaultMaxActive,
		IdleTimeout: pool.DefaultIdleTimeout,
		Wait:        true,
		WaitTimeOut: time.Second * 3,
		StatTime:    pool.DefaultStatTime,
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