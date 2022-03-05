package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"vtool/server/common"
)

func NewRegister(cluster []string) (*Register, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cluster,
		DialTimeout: common.DefaultTTl,
	})

	if err != nil {
		return nil, err
	}

	register := &Register{client: client}
	DefaultEtcdInstance = register
	return DefaultEtcdInstance, nil
}
