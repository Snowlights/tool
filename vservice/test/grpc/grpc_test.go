package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"testing"
	"time"
	"vtool/idl/grpc/grpcError"
	clientCommon "vtool/vservice/client/common"
	clientGrpc "vtool/vservice/client/grpc"
	"vtool/vservice/common"
	"vtool/vservice/server"
	. "vtool/vservice/test/grpc/grpc_protocol"
)

type helloServiceHandler struct {
}

func (h *helloServiceHandler) SayHello(ctx context.Context, req *SayHelloReq) (*SayHelloRes, error) {
	return &SayHelloRes{
		Data: &SayHelloData{Val: "this is val"},
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

var grpcClient common.RpcClient

func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(TestServiceClient) error) error {
	return grpcClient.Rpc(&common.ClientCallerArgs{
		Lane:    "",
		HashKey: hashKey,
		TimeOut: timeout,
	}, func(c interface{}) error {
		ct, ok := c.(TestServiceClient)
		if ok {
			return fn(ct)
		} else {
			return fmt.Errorf("reflect client grpc error")
		}
	})
}

func SayHello(ctx context.Context, req *SayHelloReq) (res *SayHelloRes) {
	err := rpc(ctx, "", time.Millisecond*3000,
		func(c TestServiceClient) (e error) {
			res, e = c.SayHello(ctx, req)
			return e
		})

	if err != nil {
		res = &SayHelloRes{
			ErrInfo: &grpcError.ErrInfo{
				Code: -1,
				Msg:  fmt.Sprintf("rpc service:%s serv:%s method:SayHello err:%v", "censor", common.Grpc, err),
			},
		}
	}
	return
}

func TestNewGrpcClient(t *testing.T) {

	client, _ := clientCommon.NewClientWithClientConfig(&common.ClientConfig{
		RegistrationType: common.ZOOKEEPER,
		Cluster:          []string{"127.0.0.1:2379"},
		ServGroup:        "base/talent",
		ServName:         "censor",
	})

	servCli := func(conn *grpc.ClientConn) interface{} {
		return NewTestServiceClient(conn)
	}

	grpcClient = clientGrpc.NewGrpcClient(client, servCli)
	for i := 0; i < 5000; i++ {
		go func() {
			fmt.Println(SayHello(context.Background(), &SayHelloReq{}))
		}()
	}

	time.Sleep(time.Hour)
}
