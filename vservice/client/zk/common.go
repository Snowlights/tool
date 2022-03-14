package zk

import "time"

type ClientConfig struct {
	Cluster []string
	TimeOut time.Duration

	ServGroup string
	ServName  string
}
