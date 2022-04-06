package thrift

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"testing"
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
		Data: &SayHelloData{Val: "this is thrift val"},
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
