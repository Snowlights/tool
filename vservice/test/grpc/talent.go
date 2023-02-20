package grpc

import (
	"context"
	"fmt"
	"github.com/Snowlights/tool/idl/grpc/grpcBase"
	clientCommon "github.com/Snowlights/tool/vservice/client/common"
	"github.com/Snowlights/tool/vservice/client/rpc_client"
	"github.com/Snowlights/tool/vservice/common"
	. "github.com/Snowlights/tool/vservice/test/grpc/grpc_protocol"
	"google.golang.org/grpc"
	"time"
)

func init() {
	client, _ := clientCommon.NewClientWithClientConfig(&common.ClientConfig{
		RegistrationType: common.ETCD,
		Cluster:          []string{"127.0.0.1:2379"},
		ServGroup:        "base/talent",
		ServName:         "talent",
	})

	servCli := func(conn *grpc.ClientConn) interface{} {
		return NewTestServiceClient(conn)
	}
	talentGrpcClient = rpc_client.NewRpcClient(client, nil, servCli)
}

var talentGrpcClient common.RpcClient

func talentRpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(TestServiceClient) error) error {
	return talentGrpcClient.Rpc(ctx, &common.ClientCallerArgs{
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

func TalentSayHello(ctx context.Context, req *SayHelloReq) (res *SayHelloRes) {
	err := talentRpc(ctx, "", time.Millisecond*3000,
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
