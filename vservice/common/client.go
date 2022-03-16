package common

import (
	"net/http"
	"time"
)

const (
	HashKey = "-"

	DefaultMaxTimeOut = time.Second * 10
)

var DefaultHttpClient *http.Client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 128,
		MaxConnsPerHost:     1024,
	},
	Timeout: 0,
}

type Client interface {
	GetAllServAddr() []*RegisterServiceInfo
	GetServAddr(lane string, serviceType ServiceType, hashKey string) (*ServiceInfo, bool)
}

type Caller interface {
	Do(*ClientCallerArgs, interface{}) (interface{}, error)
}

type ClientCallerArgs struct {
	Lane       string      `json:"lane"`
	ServType   ServiceType `json:"serv_type"`
	EngineType EngineType  `json:"engine_type"`
	HashKey    string      `json:"hash_key"`
}

type HttpCallerOptions struct {
	Method   string        `json:"method"`
	API      string        `json:"api"`
	Body     []byte        `json:"body"`
	Duration time.Duration `json:"duration"`
}

type ClientConfig struct {
	RegistrationType RegistrationType
	Cluster          []string

	ServGroup string
	ServName  string
}
