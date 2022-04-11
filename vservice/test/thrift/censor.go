package thrift

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
	"vtool/idl/thrift/gen-go/thriftBase"
	clientCommon "vtool/vservice/client/common"
	clientThrift "vtool/vservice/client/rpc_client"
	"vtool/vservice/common"
	. "vtool/vservice/test/thrift/thrift_protocol/gen-go/testService"
)

func init() {
	client, _ := clientCommon.NewClientWithClientConfig(&common.ClientConfig{
		RegistrationType: common.ETCD,
		Cluster:          []string{"127.0.0.1:2379"},
		ServGroup:        "base/talent",
		ServName:         "censor",
	})

	servCli := func(t thrift.TTransport, tp thrift.TProtocolFactory) interface{} {
		return NewTestServiceClientFactory(t, tp)
	}

	thriftClient = clientThrift.NewRpcClient(client, servCli, nil)
}

var thriftClient common.RpcClient

func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(ctx context.Context, c *TestServiceClient) error) error {
	return thriftClient.Rpc(ctx, &common.ClientCallerArgs{
		Lane:    "",
		HashKey: hashKey,
		TimeOut: timeout,
	}, func(ctx context.Context, c interface{}) error {
		ct, ok := c.(*TestServiceClient)
		if ok {
			return fn(ctx, ct)
		} else {
			return fmt.Errorf("reflect client rpc_client error")
		}
	})
}

func SayHello(ctx context.Context, req *SayHelloReq) (res *SayHelloRes) {
	err := rpc(ctx, "", time.Hour,
		func(ctx context.Context, c *TestServiceClient) (e error) {
			tCtx := clientCommon.NewThriftBaseContextFromContext(ctx)
			res, e = c.SayHello(req, tCtx)
			return e
		})

	if err != nil {
		res = &SayHelloRes{
			ErrInfo: &thriftBase.ErrInfo{
				Code: -1,
				Msg:  fmt.Sprintf("rpc service:%s serv:%s method:SayHello err:%v", "censor", common.Thrift, err),
			},
		}
	}
	return
}
