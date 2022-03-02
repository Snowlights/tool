package zk

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func NewRegister(cluster []string) (*Register, error) {
	conn, _, err := zk.Connect(cluster, time.Second*20)

	if err != nil {
		return nil, err
	}

	return &Register{conn: conn}, nil
}
