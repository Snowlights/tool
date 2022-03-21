package grpc

import (
	"context"
	"fmt"
	"testing"
	"vtool/idl/grpc/grpcError"
	"vtool/vservice/common"
	"vtool/vservice/server"
	. "vtool/vservice/test/grpc/grpc_protocol"
)

type helloServiceHandler struct {
}

func (h *helloServiceHandler) SayHello(ctx context.Context, req *SayHelloReq) (*SayHelloRes, error) {

	return &SayHelloRes{
		Data: &SayHelloData{Val: "this is val"},
		ErrInfo: &grpcError.ErrInfo{
			Code: -1,
			Msg:  "this is error",
		},
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
