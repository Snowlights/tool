package engine

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/Snowlights/tool/vlog"
	"github.com/Snowlights/tool/vnet"
	"github.com/Snowlights/tool/vservice/common"
)

type GinPower struct {
	c *gin.Engine
}

func (c *GinPower) Power(ctx context.Context, addr string) (string, error) {

	listener, err := vnet.ListenServAddr(ctx, addr)
	if err != nil {
		return "", err
	}
	// todo tracing and other middleware
	// this should be init before power
	serv := &http.Server{Handler: c.c}
	go func() {
		err := serv.Serve(listener)
		if err != nil {
			vlog.PanicF(ctx, "power gin engine failed, addr is %s, error is %s", listener.Addr().String(), err.Error())
		}
	}()

	return listener.Addr().String(), nil
}

func (c *GinPower) Type() common.EngineType {
	return common.Gin
}
