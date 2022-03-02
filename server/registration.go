package server

import (
	"context"
	"vtool/server/etcd"
	"vtool/server/zk"
	"vtool/vnet"
)

func RegisterService(ctx context.Context, config *RegisterConfig) error {

	// todo make group to os.env config
	path := config.Group + config.ServName
	servAddr, err := vnet.GetServAddr(config.ServAddr)
	if err != nil {
		return err
	}

	switch config.RegistrationType {
	case ETCD:
		return etcd.DefaultEtcdInstance.Register(ctx, path, servAddr, defaultTTl)
	case ZOOKEEPER:
		return zk.DefaultZkInstance.Register(ctx, path, servAddr, defaultTTl)
	default:
		return UnSupportedRegistrationType
	}

}
