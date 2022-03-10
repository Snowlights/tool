package common

import (
	"context"
)

type Processor interface {
	Prepare() error
	Engine() (string, interface{})
}

type EnginePower interface {
	Power(context.Context, string) (string, error)
	Type() string
}

type ServerBase interface {
	Register(map[string]Processor) error

	ServName() string
	ServGroup() string
	ServInfo() map[string]*ServiceInfo

	FullServiceRegisterPath() string

	Stop()

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
