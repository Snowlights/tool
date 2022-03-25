package grpc

import "google.golang.org/grpc"

type ClientConn struct {
	serviceClient interface{}
	conn          *grpc.ClientConn
}

func (c *ClientConn) Close() error {
	return c.conn.Close()
}

func (c *ClientConn) GetConn() interface{} {
	return c.serviceClient
}
