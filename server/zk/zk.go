package zk

import (
	"github.com/samuel/go-zookeeper/zk"
	"vtool/server/common"
)

func NewRegister(cluster []string) (*Register, error) {
	conn, _, err := zk.Connect(cluster, common.DefaultTTl)

	if err != nil {
		return nil, err
	}

	register := &Register{conn: conn}
	DefaultZkInstance = register
	return DefaultZkInstance, nil
}
