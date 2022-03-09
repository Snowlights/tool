package processor

import (
	"context"
	"github.com/gin-gonic/gin"
	"vtool/server/common"
	"vtool/vlog"
)

type Processor interface {
	Engine() (string, interface{})
}

func Register(ctx context.Context, props map[common.RegistrationType]Processor) error {

	// todo do init job

	// todo do listen job

	for _, processor := range props {
		addr, engine := processor.Engine()
		switch engineIns := engine.(type) {
		case *gin.Engine:
			power := ginPower{engineIns}
			servAddr, err := power.Power(ctx, addr)
			if err != nil {
				return err
			}
			// todo do register job, and get all services
			vlog.Info(ctx, "servAddr is ", servAddr)
		}
	}

	// todo wait signal and block process

	return nil
}
