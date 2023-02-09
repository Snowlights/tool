package grpc

import (
	"context"
	"fmt"
	"github.com/Snowlights/tool/vlog"
	"github.com/Snowlights/tool/vservice/common"
	"github.com/Snowlights/tool/vservice/server"
	. "github.com/Snowlights/tool/vservice/test/grpc/grpc_protocol"
	"github.com/Snowlights/tool/vsql"
	"testing"
)

const (
	cluster = "censor"
)

type helloServiceHandler struct{}

func (h *helloServiceHandler) SayHello(ctx context.Context, req *SayHelloReq) (*SayHelloRes, error) {

	//res := TalentSayHello(ctx, req)
	//
	//return &SayHelloRes{
	//	Data: &SayHelloData{Val: "this is grpc val" + fmt.Sprintf("talent res is %+v", res)},
	//}, nil

	// res := thrift.SayHello(ctx, &testService.SayHelloReq{Val: 1})
	//
	db, err := vsql.GetDB(cluster)
	if err != nil {
		return nil, fmt.Errorf("get db failed %s", err.Error())
	}

	_, err = db.ExecContext(ctx, "insert into test_table(name) values(?)", "test")
	if err != nil {
		return nil, fmt.Errorf("get db failed %s", err.Error())
	}

	redisCli := server.GetServBase().GetRedisClient(ctx)
	cmd := redisCli.Set(ctx, "test", "test", 0)
	res, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("get db failed %s", err.Error())
	}
	vlog.ErrorF(ctx, "grpc say hello req is %+v, redis res is %+v", req, res)

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

	fmt.Println(SayHello(ctx, &SayHelloReq{}))

	// time.Sleep(time.Hour)
}
