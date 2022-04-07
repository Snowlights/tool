package thrift

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"testing"
	"vtool/idl/thrift/gen-go/thriftBase"
	common2 "vtool/vservice/client/common"
	"vtool/vservice/common"
	"vtool/vservice/server"
	"vtool/vservice/test/grpc"
	testService "vtool/vservice/test/grpc/grpc_protocol"
	. "vtool/vservice/test/thrift/thrift_protocol/gen-go/testService"
)

type helloServiceHandler struct {
}

func (h *helloServiceHandler) SayHello(req *SayHelloReq, tctx *thriftBase.Context) (*SayHelloRes, error) {
	ctx := common2.NewContextFromThriftBaseContext("helloServiceHandler.SayHello", tctx)
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		defer span.Finish()
	}

	res := grpc.SayHello(ctx, &testService.SayHelloReq{HelloType: 1})

	return &SayHelloRes{
		Data: &SayHelloData{Val: "this is thrift val" + fmt.Sprintf("%+v", res)},
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
	fmt.Println(SayHello(context.Background(), &SayHelloReq{}))
}
