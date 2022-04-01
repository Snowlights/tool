package thrift

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"testing"
	"time"
	"vtool/idl/thrift/gen-go/thriftError"
	clientCommon "vtool/vservice/client/common"
	clientThrift "vtool/vservice/client/rpc_client"
	"vtool/vservice/common"
	"vtool/vservice/server"
	. "vtool/vservice/test/thrift/thrift_protocol/gen-go/testService"
)

type helloServiceHandler struct {
}

func (h *helloServiceHandler) SayHello(req *SayHelloReq) (*SayHelloRes, error) {
	return &SayHelloRes{
		Data: &SayHelloData{Val: "this is val"},
	}, nil
}

type TestProcessor struct{}

func (mp *TestProcessor) Engine() (string, interface{}) {

	return "", NewTestServiceProcessor(new(helloServiceHandler))
}

func TestThriftServer(t *testing.T) {
	err := server.ServService(map[common.ServiceType]common.Processor{
		common.Thrift: &TestProcessor{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

}

var thriftClient common.RpcClient

func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(*TestServiceClient) error) error {
	return thriftClient.Rpc(ctx, &common.ClientCallerArgs{
		Lane:    "",
		HashKey: hashKey,
		TimeOut: timeout,
	}, func(c interface{}) error {
		ct, ok := c.(*TestServiceClient)
		if ok {
			return fn(ct)
		} else {
			return fmt.Errorf("reflect client rpc_client error")
		}
	})
}

func SayHello(ctx context.Context, req *SayHelloReq) (res *SayHelloRes) {
	err := rpc(ctx, "", time.Millisecond*3000,
		func(c *TestServiceClient) (e error) {
			res, e = c.SayHello(req)
			return e
		})

	if err != nil {
		res = &SayHelloRes{
			ErrInfo: &thriftError.ErrInfo{
				Code: -1,
				Msg:  fmt.Sprintf("rpc service:%s serv:%s method:SayHello err:%v", "censor", common.Thrift, err),
			},
		}
	}
	return
}

func TestNewThriftClient(t *testing.T) {

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

	fmt.Println(SayHello(context.Background(), &SayHelloReq{}))
}
