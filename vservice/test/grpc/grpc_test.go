package grpc

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"testing"
	"vtool/vservice/common"
	"vtool/vservice/server"
	. "vtool/vservice/test/grpc/grpc_protocol"
)

type helloServiceHandler struct{}

func (h *helloServiceHandler) SayHello(ctx context.Context, req *SayHelloReq) (*SayHelloRes, error) {

	//res := TalentSayHello(ctx, req)
	//
	//return &SayHelloRes{
	//	Data: &SayHelloData{Val: "this is grpc val" + fmt.Sprintf("talent res is %+v", res)},
	//}, nil

	return &SayHelloRes{
		Data: &SayHelloData{Val: "this is grpc val"},
	}, nil
}

type TestProcessor struct{}

func (mp *TestProcessor) Engine() (string, interface{}) {
	serv := server.NewGrpcServerWithInterceptor()
	RegisterTestServiceServer(serv, new(helloServiceHandler))
	return "", serv
}

func TestGrpcServer(t *testing.T) {

	err := server.ServService(map[common.ServiceType]common.Processor{
		common.Grpc: &TestProcessor{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestGrpcServer2(t *testing.T) {

	err := server.ServService(map[common.ServiceType]common.Processor{
		common.Grpc: &TestProcessor{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestNewGrpcClient(t *testing.T) {
	ctx := context.Background()
	opentracing.SpanFromContext(ctx)

	fmt.Println(SayHello(ctx, &SayHelloReq{}))

	// time.Sleep(time.Hour)
}
