package common

import (
	"github.com/Snowlights/tool/vservice/client/etcd"
	"github.com/Snowlights/tool/vservice/client/zk"
	"github.com/Snowlights/tool/vservice/common"
	"math/rand"
	"strconv"
	"time"
)

func NewClientWithClientConfig(cliConfig *common.ClientConfig) (common.Client, error) {
	switch cliConfig.RegistrationType {
	case common.ETCD:
		client, err := etcd.NewEtcdClient(&etcd.ClientConfig{
			Cluster:   cliConfig.Cluster,
			TimeOut:   common.DefaultTTl,
			ServGroup: cliConfig.ServGroup,
			ServName:  cliConfig.ServName,
		})
		if err != nil {
			return nil, err
		}
		return client, nil
	case common.ZOOKEEPER:
		client, err := zk.NewZkClient(&zk.ClientConfig{
			Cluster:   cliConfig.Cluster,
			TimeOut:   common.DefaultTTl,
			ServGroup: cliConfig.ServGroup,
			ServName:  cliConfig.ServName,
		})
		if err != nil {
			return nil, err
		}
		return client, nil
	default:
		return nil, common.UnSupportedRegistrationType
	}
}

func NewHashKey() string {
	rand.Seed(int64(time.Now().Nanosecond()))
	return strconv.FormatInt(rand.Int63(), 10)
}
