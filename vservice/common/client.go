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
	AddPoolHandler(func([]string))
}

type RpcConn interface {
	GetConn() interface{}
	Close() error
}

type ClientCallerArgs struct {
	Lane    string `json:"lane"`
	HashKey string `json:"hash_key"`
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
