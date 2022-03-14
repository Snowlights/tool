package zk

import (
	"time"
)

type RegisterConfig struct {
	Cluster []string
	TimeOut time.Duration
}

// must init before app init todo addr change to apollo config
