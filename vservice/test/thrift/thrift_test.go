package thrift

import (
	"fmt"
	"testing"
	"vtool/idl/thrift/gen-go/thriftError"
	"vtool/vservice/common"
	"vtool/vservice/server"
	. "vtool/vservice/test/thrift/thrift_protocol/gen-go/testService"
)

type helloServiceHandler struct {
}

func (h *helloServiceHandler) SayHello(req *SayHelloReq) (*SayHelloRes, error) {
	return &SayHelloRes{
		Data: &SayHelloData{Val: "this is val"},
		ErrInfo: &thriftError.ErrInfo{
			Code: 1,
			Msg:  "this is error",
		},
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
