package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
	"vtool/idl/grpc/grpcBase"
	clientCommon "vtool/vservice/client/common"
	"vtool/vservice/client/rpc_client"
	"vtool/vservice/common"
	. "vtool/vservice/test/grpc/grpc_protocol"
)

func init() {
	client, _ := clientCommon.NewClientWithClientConfig(&common.ClientConfig{
		RegistrationType: common.ETCD,
		Cluster:          []string{"127.0.0.1:2379"},
		ServGroup:        "base/talent",
		ServName:         "censor",
	})

	servCli := func(conn *grpc.ClientConn) interface{} {
		return NewTestServiceClient(conn)
	}
	censorGrpcClient = rpc_client.NewRpcClient(client, nil, servCli)
}

var censorGrpcClient common.RpcClient

func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(TestServiceClient) error) error {
	return censorGrpcClient.Rpc(ctx, &common.ClientCallerArgs{
		Lane:    "",
		HashKey: hashKey,
		TimeOut: timeout,
	}, func(ctx context.Context, c interface{}) error {
		ct, ok := c.(TestServiceClient)
		if ok {
			return fn(ct)
		} else {
			return fmt.Errorf("reflect client grpc error")
		}
	})
}

func SayHello(ctx context.Context, req *SayHelloReq) (res *SayHelloRes) {
	err := rpc(ctx, "", time.Hour,
		func(c TestServiceClient) (e error) {
			res, e = c.SayHello(ctx, req)
			return e
		})

	if err != nil {
		res = &SayHelloRes{
			ErrInfo: &grpcBase.ErrInfo{
				Code: -1,
				Msg:  fmt.Sprintf("rpc service:%s serv:%s method:SayHello err:%v", "censor", common.Grpc, err),
			},
		}
	}
	return
}
