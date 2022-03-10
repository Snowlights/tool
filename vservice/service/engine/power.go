package engine

import (
	"github.com/gin-gonic/gin"
	"vtool/vservice/common"
)

func GetEnginePower(engine interface{}) (common.EnginePower, bool) {

	switch engineIns := engine.(type) {
	case *gin.Engine:
		enginePower := &GinPower{engineIns}
		return enginePower, true
	}

	return nil, false
}
