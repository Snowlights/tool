package client

import (
	"time"
	"vtool/vservice/client/etcd"
	"vtool/vservice/client/zk"
	"vtool/vservice/common"
)

// todo client lookup

type HttpClient struct {
	registerType common.RegistrationType
	serviceType  common.ServiceType

	client common.Client
}

func NewHttpClient(cliConfig *common.ClientConfig) (common.Client, error) {
	switch cliConfig.RegistrationType {
	case common.ETCD:
		return etcd.NewEtcdClient(&etcd.ClientConfig{
			Cluster:   cliConfig.Cluster,
			TimeOut:   common.DefaultTTl,
			ServGroup: cliConfig.ServGroup,
			ServName:  cliConfig.ServName,
		})
	case common.ZOOKEEPER:
		return zk.NewZkClient(&zk.ClientConfig{
			Cluster:   cliConfig.Cluster,
			TimeOut:   common.DefaultTTl,
			ServGroup: cliConfig.ServGroup,
			ServName:  cliConfig.ServName,
		})
	default:
		return nil, common.UnSupportedRegistrationType
	}

}

func (hc HttpClient) Do(engineType common.EngineType, runFunc func(addr string, timeout time.Duration)) {
	// todo find a ins to call runFunc
}
