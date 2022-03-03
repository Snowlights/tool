package vcpu

import (
	"context"
	"vtool/vlog"
)

type CPU struct {
	Cores int32
	Mhz   float64
}

type CPUMonitor interface {
	Usage() (float64, error)
	Info() (*CPU, error)
}

const (
	defaultInterval = 0
)

var defaultCPUInstance CPUMonitor

func init() {
	ins, err := NewPsutilCPU(defaultInterval)
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	defaultCPUInstance = ins
}
