package thrift

import "git.apache.org/thrift.git/lib/go/thrift"

type ClientConn struct {
	thriftSocket  *thrift.TSocket
	transport     thrift.TTransport
	serviceClient interface{}
}

func (tcc *ClientConn) Close() error {
	return tcc.transport.Close()
}

func (tcc *ClientConn) GetConn() interface{} {
	return tcc.serviceClient
}
