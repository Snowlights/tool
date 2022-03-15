package register

import (
	"vtool/vservice/common"
	"vtool/vservice/server/register/etcd"
	"vtool/vservice/server/register/zk"
)

func GetRegisterEngine(registerConfig *common.RegisterConfig) (common.Register, error) {
	switch registerConfig.RegistrationType {
	case common.ETCD:
		return etcd.NewRegister(&etcd.RegisterConfig{
			Cluster: registerConfig.Cluster,
			TimeOut: common.DefaultTTl,
		})
	case common.ZOOKEEPER:
		return zk.NewRegister(&zk.RegisterConfig{
			Cluster: registerConfig.Cluster,
			TimeOut: common.DefaultTTl,
		})
	default:
		return nil, common.UnSupportedRegistrationType
	}

}
