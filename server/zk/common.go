package zk

import (
	"context"
	"vtool/vlog"
)

const (
	character = ":"
	Slash     = "/"
)

var defaultCluster = []string{"127.0.0.1:2181"}

var DefaultZkInstance *Register

func init() {
	ins, err := NewRegister(defaultCluster)
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	DefaultZkInstance = ins
}