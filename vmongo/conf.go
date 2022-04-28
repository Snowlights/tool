package vmongo

import (
	"time"
)

const (
	defaultTimeout      = time.Second * 5
	defaultReadTimeout  = time.Second * 5
	defaultWriteTimeout = time.Second * 5
	defaultPoolLimit    = 128
)
