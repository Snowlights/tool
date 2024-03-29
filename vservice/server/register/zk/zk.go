package zk

import (
	"github.com/Snowlights/tool/vservice/common"
	"github.com/samuel/go-zookeeper/zk"
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
