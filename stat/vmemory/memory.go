package vmemory

func NewMemory() (MemoryMonitor, error) {
	memory := NewRuntimeMemory()
	defaultMemoryInstance = memory
	return defaultMemoryInstance, nil
}

func Virtual() uint64 {
	return defaultMemoryInstance.Virtual()
}

func GoroutineNums() int64 {
	return defaultMemoryInstance.GoroutineNums()
}

func HeapAlloc() int64 {
	return defaultMemoryInstance.HeapAlloc()
}

func HeapObjects() int64 {
	return defaultMemoryInstance.HeapObjects()
}

func GCPause() float64 {
	return defaultMemoryInstance.GCPause()
}
