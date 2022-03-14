package common

import "context"

type RegisterServiceInfo struct {
	ServPath string
	ServList map[ServiceType]*ServiceInfo
}

type ServiceInfo struct {
	Type EngineType `json:"type"`
	Addr string     `json:"addr"`
}

type ServerBase interface {
	Register(context.Context, map[ServiceType]Processor) error

	ServName() string
	ServGroup() string
	ServInfo() map[ServiceType]*ServiceInfo

	FullServiceRegisterPath() string

	Stop()

	// todo add config to apollo

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
