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
