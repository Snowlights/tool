package etcd

import (
	"context"
	"errors"
	"vtool/vlog"
)

const (
	leaseSuccess = "Lease renewal succeeded"
	character    = ":"
	equals       = "="
)

var lockFailed = errors.New("lock failed")

var defaultCluster = []string{"127.0.0.1:2379"}

var DefaultEtcdInstance *Register

func init() {
	ins, err := NewRegister(defaultCluster)
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	DefaultEtcdInstance = ins
}
