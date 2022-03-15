package client

import (
	"net/http"
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

func (hc HttpClient) Do(engineType common.EngineType) (*http.Response, error) {
	// todo find a ins to call runFunc

	//serv, ok  := hc.client.GetServAddr("", "", "")
	//if !ok {
	//	return nil, nil
	//}

	return nil, nil
}

func (hc HttpClient) do(serv *common.ServiceInfo) {

	// todo  do http request with context timeout

}
