package common

type Client interface {
	GetAllServAddr() []*RegisterServiceInfo
}

type ClientConfig struct {
	RegistrationType RegistrationType
	Cluster          []string

	ServGroup string
	ServName  string
}
