package etcd

import (
	"context"
	"errors"
	"time"
	"vtool/server/common"
	"vtool/vlog"
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
