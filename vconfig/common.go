package vconfig

import (
	"github.com/apolloconfig/agollo/v4/storage"
	"time"
)

const (
	// server config
	Application = "application"
	// client config
	Client = "client"
	// db config
	DB = "db"
	// mq config
	Mq = "mq"
	// log config
	Log = "log"
)

type Center interface {
	GetValue(key string) string
	GetValueWithNamespace(namespace, key string) string
	UnmarshalWithNameSpace(namespace, tag string, v interface{}) error
	AddListener(listeners storage.ChangeListener)
}

type ClientConfig struct {
	Idle        int64         `json:"idle" properties:"idle"`
	IdleTimeout int64         `json:"idle_timeout" properties:"idle_timeout"`
	MaxActive   time.Duration `json:"max_active" properties:"max_active"`
	StatTime    int64         `json:"stat_time" properties:"stat_time"`
	Wait        bool          `json:"wait" properties:"wait"`
	WaitTimeout int64         `json:"wait_timeout" properties:"wait_timeout"`
}

type LogConfig struct {
	Level string `json:"level" properties:"level"`
	Path  string `json:"path" properties:"path"`
}
