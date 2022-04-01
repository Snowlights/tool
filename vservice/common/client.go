package common

import (
	"context"
	"github.com/apolloconfig/agollo/v4/storage"
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
	ServName() string
	ServGroup() string
}

type RpcClient interface {
	Rpc(context.Context, *ClientCallerArgs, func(interface{}) error) error
}

type RpcConn interface {
	GetConn() interface{}
	Close() error
}

type ClientCallerArgs struct {
	Lane    string        `json:"lane"`
	HashKey string        `json:"hash_key"`
	TimeOut time.Duration `json:"time_out"`
}

type HttpCallerOptions struct {
	Method string `json:"method"`
	API    string `json:"api"`
	Body   []byte `json:"body"`
}

type ClientConfig struct {
	RegistrationType RegistrationType
	Cluster          []string

	ServGroup string
	ServName  string
}

type ClientListener struct {
	Change func()
}

func (cl *ClientListener) OnChange(event *storage.ChangeEvent) {

}

func (cl *ClientListener) OnNewestChange(event *storage.FullChangeEvent) {
	cl.Change()
}
