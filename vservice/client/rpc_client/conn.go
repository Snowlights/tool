package rpc_client

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"google.golang.org/grpc"
)

type ThriftClientConn struct {
	thriftSocket  *thrift.TSocket
	transport     thrift.TTransport
	serviceClient interface{}
}

func (c *ThriftClientConn) Close() error {
	return c.transport.Close()
}

func (c *ThriftClientConn) GetConn() interface{} {
	return c.serviceClient
}

type GrpcClientConn struct {
	serviceClient interface{}
	conn          *grpc.ClientConn
}

func (c *GrpcClientConn) Close() error {
	return c.conn.Close()
}

func (c *GrpcClientConn) GetConn() interface{} {
	return c.serviceClient
}
