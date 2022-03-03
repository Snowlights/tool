package vcpu

import "time"

func NewCPUMonitor(interval time.Duration) (CPUMonitor, error) {
	return NewPsutilCPU(interval)
}

func Usage() (float64, error) {
	return defaultCPUInstance.Usage()
}

func Info() (*CPU, error) {
	return defaultCPUInstance.Info()
}
