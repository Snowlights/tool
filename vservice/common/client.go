package common

import "time"

const (
	HashKey = "-"

	DefaultMaxTimeOut = time.Second * 10
)

type Client interface {
	GetAllServAddr() []*RegisterServiceInfo
	GetServAddr(lane string, serviceType ServiceType, hashKey string) (*ServiceInfo, bool)
}

type ClientConfig struct {
	RegistrationType RegistrationType
	Cluster          []string

	ServGroup string
	ServName  string
}
