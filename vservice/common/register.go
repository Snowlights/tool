package common

import (
	"context"
	"time"
)

type (
	RegistrationType int64
	ServiceType      string
)

const (
	Colon      = ":"
	Equals     = "="
	Slash      = "/"
	Bar        = "-"
	HttpPrefix = "http://"

	DefaultRegisterPath = "/tools"
)

const (
	ETCD      RegistrationType = 1
	ZOOKEEPER RegistrationType = 2
	Consul    RegistrationType = 3

	HTTP   ServiceType = "http"
	Thrift ServiceType = "thrift"
	Grpc   ServiceType = "grpc"

	Metric ServiceType = "metric"
)

const DefaultTTl = time.Second * 10

type RegisterConfig struct {
	RegistrationType RegistrationType
	Cluster          []string
}

type Register interface {
	// get key val
	Get(ctx context.Context, path string) (string, error)
	// Refresh the expiration time of the node without updating the value
	RefreshTtl(ctx context.Context, path, val string, ttl time.Duration) error
	// Execute registration, and the heartbeat will be maintained after registration.
	// When calling registration, the current value will be set to the node
	Register(ctx context.Context, path, val string, ttl time.Duration) (string, error)
	// unRegister Service
	UnRegister(ctx context.Context, path string) error
	// get all node
	GetNode(ctx context.Context, path string) ([]*RegisterServiceInfo, error)
}
