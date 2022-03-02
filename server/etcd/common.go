package etcd

import (
	"context"
	"vtool/vlog"
)

const (
	leaseSuccess = "续租成功"
	character    = ":"
)

var defaultCluster = []string{"127.0.0.1:2379"}

var DefaultEtcdInstance *Register

func init() {
	ins, err := NewRegister(defaultCluster)
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	DefaultEtcdInstance = ins
}
