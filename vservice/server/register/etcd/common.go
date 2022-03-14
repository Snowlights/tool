package etcd

import (
	"context"
	"errors"
	"time"
	"vtool/vlog"
	"vtool/vservice/common"
)

const (
	leaseSuccess = "Lease renewal succeeded"
)

type RegisterConfig struct {
	Cluster []string
	TimeOut time.Duration
}

var lockFailed = errors.New("lock failed")

var defaultCluster = []string{"127.0.0.1:2379"}

var DefaultEtcdInstance *Register

// must init before app init todo addr change to apollo config

func init() {
	ins, err := NewRegister(&RegisterConfig{
		Cluster: defaultCluster,
		TimeOut: common.DefaultTTl,
	})
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	DefaultEtcdInstance = ins
}
