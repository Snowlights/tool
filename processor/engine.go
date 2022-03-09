package processor

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"vtool/vlog"
	"vtool/vnet"
)

type Processor interface {
	Prepare() error
	Engine() (string, interface{})
}

type EnginePower interface {
	Power(context.Context, string) error
}

type ginPower struct {
	c *gin.Engine
}

func (c *ginPower) Power(ctx context.Context, addr string) (string, error) {

	listener, err := vnet.ListenServAddr(ctx, addr)
	if err != nil {
		return "", err
	}
	// todo tracing and other middleware
	serv := &http.Server{Handler: c.c}
	go func() {
		err := serv.Serve(listener)
		if err != nil {
			vlog.PanicF(ctx, "power gin engine failed, addr is %s, error is %s", listener.Addr().String(), err.Error())
		}
	}()

	return listener.Addr().String(), nil
}
