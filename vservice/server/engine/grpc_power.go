package engine

import (
	"context"
	"github.com/Snowlights/tool/vlog"
	"github.com/Snowlights/tool/vnet"
	"github.com/Snowlights/tool/vservice/common"
	"google.golang.org/grpc"
	"time"
)

type GrpcPower struct {
	c *grpc.Server
}

func (c *GrpcPower) Power(ctx context.Context, addr string) (string, error) {

	listener, err := vnet.ListenServAddr(ctx, addr)
	if err != nil {
		return "", err
	}

	go func() {
		err := c.c.Serve(listener)
		if err != nil {
			vlog.PanicF(ctx, "GrpcPower.Power failed, error is %s, addr is %+v", err.Error(), listener.Addr().String())
		}
	}()

	return listener.Addr().String(), nil
}

func newDisableContextCancelGrpcUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		valueCtx := NewValueContext(ctx)
		return handler(valueCtx, req)
	}
}

type ValueContext struct {
	ctx context.Context
}

func (c ValueContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c ValueContext) Done() <-chan struct{}             { return nil }
func (c ValueContext) Err() error                        { return nil }
func (c ValueContext) Value(key interface{}) interface{} { return c.ctx.Value(key) }

// NewValueContext returns a context that is never canceled.
func NewValueContext(ctx context.Context) context.Context {
	return ValueContext{ctx: ctx}
}

func (c *GrpcPower) Type() common.EngineType {
	return common.Rpc
}
