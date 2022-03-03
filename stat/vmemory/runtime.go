package vmemory

import "runtime"

type RuntimeMemory struct {
	stat runtime.MemStats
}

func NewRuntimeMemory() *RuntimeMemory {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	return &RuntimeMemory{stat: stat}
}

func (r *RuntimeMemory) Virtual() uint64 {
	return r.stat.Sys
}

func (r *RuntimeMemory) GoroutineNums() int64 {
	return int64(runtime.NumGoroutine())
}

func (r *RuntimeMemory) HeapAlloc() int64 {
	return int64(r.stat.HeapAlloc)
}

func (r *RuntimeMemory) HeapObjects() int64 {
	return int64(r.stat.HeapObjects)
}

func (r *RuntimeMemory) GCPause() float64 {
	return float64(r.stat.PauseNs[(r.stat.NumGC-1)%uint32(len(r.stat.PauseNs))])
}
