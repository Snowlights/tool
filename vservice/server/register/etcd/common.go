package etcd

import (
	"errors"
	"time"
)

const (
	leaseSuccess = "Lease renewal succeeded"
)

type RegisterConfig struct {
	Cluster []string
	TimeOut time.Duration
}

var lockFailed = errors.New("lock failed")
