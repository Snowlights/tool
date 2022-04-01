package engine

import (
	"context"
	"git.apache.org/thrift.git/lib/go/thrift"
	"vtool/vlog"
	"vtool/vnet"
	"vtool/vservice/common"
)

type ThriftPower struct {
	c thrift.TProcessor
}

func (c *ThriftPower) Power(ctx context.Context, addr string) (string, error) {

	listenAddr, err := vnet.GetServAddr(addr)
	if err != nil {
		return "", err
	}
	// todo tracing and other middleware
	// this should be init before power
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	// protocolFactory := rpc_client.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(listenAddr)
	if err != nil {
		return "", err
	}

	server := thrift.NewTSimpleServer4(c.c, serverTransport, transportFactory, protocolFactory)
	err = serverTransport.Listen()
	if err != nil {
		return "", err
	}

	servAddr, err := vnet.GetServAddr(serverTransport.Addr().String())
	if err != nil {
		return "", err
	}

	go func() {
		err := server.Serve()
		if err != nil {
			vlog.PanicF(ctx, "ThriftPower.Power failed, error is %s, addr is %+v", err.Error(), servAddr)
		}
	}()

	return servAddr, nil
}

func (c *ThriftPower) Type() common.EngineType {
	return common.Rpc
}
