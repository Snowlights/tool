package etcd

import (
	"github.com/Snowlights/tool/vservice/common"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewRegister(regConfig *RegisterConfig) (*Register, error) {

	timeOut := common.DefaultTTl
	if regConfig.TimeOut > 0 {
		timeOut = regConfig.TimeOut
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   regConfig.Cluster,
		DialTimeout: timeOut,
	})

	if err != nil {
		return nil, err
	}

	register := &Register{client: client}
	return register, nil
}
