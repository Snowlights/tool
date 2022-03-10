package zk

import (
	"context"
	"time"
	"vtool/vlog"
	"vtool/vservice/common"
)

type RegisterConfig struct {
	Cluster []string
	TimeOut time.Duration
}

var defaultCluster = []string{"127.0.0.1:2181"}

var DefaultZkInstance *Register

func init() {
	ins, err := NewRegister(&RegisterConfig{
		Cluster: defaultCluster,
		TimeOut: common.DefaultTTl,
	})
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	DefaultZkInstance = ins
}
