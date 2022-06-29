package thrift

import (
	"context"
	"fmt"
	"github.com/Snowlights/tool/idl/thrift/gen-go/thriftBase"
	common2 "github.com/Snowlights/tool/vservice/client/common"
	"github.com/Snowlights/tool/vservice/common"
	"github.com/Snowlights/tool/vservice/server"
	. "github.com/Snowlights/tool/vservice/test/thrift/thrift_protocol/gen-go/testService"
	"github.com/Snowlights/tool/vsql"
	"github.com/opentracing/opentracing-go"
	"testing"
)

const (
	cluster = "censor"
)

type helloServiceHandler struct {
}

func (h *helloServiceHandler) SayHello(req *SayHelloReq, tctx *thriftBase.Context) (*SayHelloRes, error) {
	ctx := common2.NewContextFromThriftBaseContext("helloServiceHandler.SayHello", tctx)
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		defer span.Finish()
	}

	db, err := vsql.GetDB(cluster)
	if err != nil {
		return nil, fmt.Errorf("get db failed %s", err.Error())
	}

	_, err = db.ExecContext(ctx, "insert into test_table(name) values(?)", "test")
	if err != nil {
		return nil, fmt.Errorf("get db failed %s", err.Error())
	}

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

func TestNewGrpcClient(t *testing.T) {
	ctx := context.Background()

	fmt.Println(SayHello(ctx, &SayHelloReq{}))

	// time.Sleep(time.Hour)
}
