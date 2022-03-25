package thrift

import "git.apache.org/thrift.git/lib/go/thrift"

type ClientConn struct {
	thriftSocket  *thrift.TSocket
	transport     thrift.TTransport
	serviceClient interface{}
}

func (c *ClientConn) Close() error {
	return c.transport.Close()
}

func (c *ClientConn) GetConn() interface{} {
	return c.serviceClient
}
