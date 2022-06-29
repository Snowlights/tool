package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"github.com/Snowlights/tool/vservice/common"
)

func GetEnginePower(engine interface{}) (common.EnginePower, bool) {

	switch engineIns := engine.(type) {
	case *gin.Engine:
		enginePower := &GinPower{engineIns}
		return enginePower, true
	case *grpc.Server:
		enginePower := &GrpcPower{engineIns}
		return enginePower, true
	case thrift.TProcessor:
		enginePower := &ThriftPower{engineIns}
		return enginePower, true
	}

	return nil, false
}
