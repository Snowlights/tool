package common

import (
	"context"
	"vtool/vconfig"
)

const (
	ServWeight = 100

	ContentKey = "vtool.content"
)

type RegisterServiceInfo struct {
	ServPath string                       `json:"serv_path"`
	Lane     string                       `json:"lane"`
	ServList map[ServiceType]*ServiceInfo `json:"serv_list"`
}

type ServiceInfo struct {
	Type EngineType `json:"type"`
	Addr string     `json:"addr"`
}

type ServerBase interface {
	Register(context.Context, map[ServiceType]Processor) error

	ServName() string
	ServGroup() string
	ServInfo() *RegisterServiceInfo

	FullServiceRegisterPath() string

	Stop()

	GetCenter(ctx context.Context) vconfig.Center

	// todo service region
	// eg: beijing、hangzhou、shanghai

	// todo service lane
	// eg：test1 lane for function1；test2 lane for function2

	// todo same service lock
	// for service selection
	// eg: master-slave model

	// todo global service lock
	// for something I do not know now

}
