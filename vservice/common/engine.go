package common

import (
	"context"
)

type EngineType string

const (
	Gin EngineType = "gin"
)

type Processor interface {
	Engine() (string, interface{})
}

type EnginePower interface {
	Power(context.Context, string) (string, error)
	Type() EngineType
}
