package vmemory

import (
	"context"
	"github.com/Snowlights/tool/vlog"
)

// calculate current process

type MemoryMonitor interface {
	Virtual() uint64
	GoroutineNums() int64
	HeapAlloc() int64
	HeapObjects() int64
	GCPause() float64
}

var defaultMemoryInstance MemoryMonitor

func init() {
	ins, err := NewMemory()
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	defaultMemoryInstance = ins
}
