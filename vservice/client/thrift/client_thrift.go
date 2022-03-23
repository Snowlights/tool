package thrift

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"vtool/vlog"
	"vtool/vservice/client/pool"
	"vtool/vservice/common"
)

// todo thrift client

type ThriftClient struct {
	client common.Client

	serviceClient func(thrift.TTransport, thrift.TProtocolFactory) interface{}

	clientPool *pool.ClientPool
}

func NewThriftClient(client common.Client, servCli func(thrift.TTransport, thrift.TProtocolFactory) interface{}) *ThriftClient {

	tc := &ThriftClient{
		client:        client,
		serviceClient: servCli,
	}

	tc.client.AddPoolHandler(tc.deleteAddrHandler)
	return tc
}

func (t *ThriftClient) deleteAddrHandler(addr []string) {
	for _, addr := range addr {
		t.clientPool.Delete(context.Background(), addr)
	}
}

func (t *ThriftClient) Rpc(args *common.ClientCallerArgs, fnRpc func(interface{}) error) (interface{}, error) {

	serv, ok := t.client.GetServAddr(args.Lane, common.Thrift, args.HashKey)
	if !ok {
		return nil, fmt.Errorf("%s caller args is %+v", common.NotFoundServInfo, args)
	}
	if serv.Type != common.Rpc {
		return nil, fmt.Errorf("%s serv info is %+v, caller args is %+v", common.NotFoundServEngine, serv, args)
	}

	return nil, t.do(context.TODO(), serv, fnRpc)
}

func (t *ThriftClient) do(ctx context.Context, serv *common.ServiceInfo, fnRpc func(interface{}) error) error {

	var err error
	retry := 3
	for ; retry >= 0; retry-- {
		err = t.rpc(ctx, serv, fnRpc)
		if err == nil {
			return nil
		}
	}
	return err
}

func (t *ThriftClient) rpc(ctx context.Context, serv *common.ServiceInfo, fnRpc func(interface{}) error) error {
	conn, err := t.clientPool.Get(ctx, serv)
	if err != nil {
		return err
	}

	return fnRpc(conn.GetConn())
}

func (t *ThriftClient) newConn(addr string) (common.RpcConn, error) {
	ctx := context.Background()

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(addr)
	if err != nil {
		vlog.ErrorF(ctx, "open thriftSocket addr:%s, err:%v", addr, err)
		return nil, err
	}
	useTransport := transportFactory.GetTransport(transport)

	if err := useTransport.Open(); err != nil {
		vlog.ErrorF(ctx, "open addr:%s err:%v", addr, err)
		return nil, err
	}

	return &ClientConn{
		thriftSocket:  transport,
		transport:     useTransport,
		serviceClient: t.serviceClient(useTransport, protocolFactory),
	}, nil
}
