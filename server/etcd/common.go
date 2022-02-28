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

var DefaulEtcdInstance *Client

func init() {
	ins, err := NewClient(defaultCluster)
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	DefaulEtcdInstance = ins
}