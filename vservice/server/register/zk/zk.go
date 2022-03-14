package zk

import (
	"github.com/samuel/go-zookeeper/zk"
	"vtool/vservice/common"
)

func NewRegister(regConfig *RegisterConfig) (*Register, error) {
	timeOut := common.DefaultTTl
	if regConfig.TimeOut > 0 {
		timeOut = regConfig.TimeOut
	}
	conn, _, err := zk.Connect(regConfig.Cluster, timeOut)

	if err != nil {
		return nil, err
	}

	register := &Register{conn: conn}
	return register, nil
}
