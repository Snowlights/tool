package vcpu

import (
	"github.com/shirou/gopsutil/cpu"
	"time"
)

type PsutilCPU struct {
	interval time.Duration
}

func NewPsutilCPU(interval time.Duration) (CPUMonitor, error) {
	psutilCPU := &PsutilCPU{interval: interval}
	defaultPsutilCPUInstance = psutilCPU
	return defaultPsutilCPUInstance, nil
}

func (p *PsutilCPU) Usage() (float64, error) {

	percent, err := cpu.Percent(p.interval, true)
	if err != nil {
		return 0, err
	}

	return percent[0], nil
}

func (p *PsutilCPU) Info() (*CPU, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	return &CPU{
		Cores: info[0].Cores,
		Mhz:   info[0].Mhz,
	}, nil
}
