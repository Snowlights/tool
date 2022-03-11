package register

import (
	"vtool/vservice/common"
	"vtool/vservice/server/register/consul"
	"vtool/vservice/server/register/etcd"
	"vtool/vservice/server/register/zk"
)

func GetRegisterEngine(registerType common.RegistrationType) (common.Register, error) {

	var engine common.Register
	switch registerType {
	case common.ETCD:
		engine = etcd.DefaultEtcdInstance
	case common.ZOOKEEPER:
		engine = zk.DefaultZkInstance
	case common.Consul:
		// only for metric collection
		engine = consul.DefaultConsulInstance
	default:
		return nil, common.UnSupportedRegistrationType
	}

	return engine, nil
}
