package server

import (
	"context"
	"time"
)

type RegistrationType int64

const (
	ETCD      RegistrationType = 1
	ZOOKEEPER RegistrationType = 2
)

const (
	defaultTTl = time.Second * 20

	_defaultID = "-1"
	retryTime  = 4
)

type RegisterConfig struct {
	RegistrationType RegistrationType

	ServName string
	ServAddr string

	Group string
}

type Register interface {
	// get key val
	Get(ctx context.Context, path string) (string, error)
	// Refresh the expiration time of the node without updating the value
	RefreshTtl(ctx context.Context, path, val string, ttl time.Duration) error
	// Execute registration, and the heartbeat will be maintained after registration.
	// When calling registration, the current value will be set to the node
	Register(ctx context.Context, path, val string, ttl time.Duration) error
}
