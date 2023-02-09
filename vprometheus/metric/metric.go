package metric

import (
	"context"
	"github.com/Snowlights/tool/stat/vcpu"
	"github.com/Snowlights/tool/stat/vload"
	"github.com/Snowlights/tool/stat/vmemory"
	"time"
)

const (
	reloadTime = time.Second * 10
)

var defaultProcessMonitor *ProcessMonitor

func init() {
	defaultProcessMonitor = &ProcessMonitor{}
}

type ProcessMonitor struct {
	Group    string
	Service  string
	Instance string
}

func InitBaseMetric(ctx context.Context, group, service, instance string) *ProcessMonitor {
	p := &ProcessMonitor{
		Group:    group,
		Service:  service,
		Instance: instance,
	}
	defaultProcessMonitor = p
	go p.reload(ctx)
	return p
}

func (p ProcessMonitor) reload(ctx context.Context) {
	for {
		p.work()
		select {
		case <-time.After(reloadTime):
		case <-ctx.Done():
			break
		}
	}
}

func (p ProcessMonitor) statApi(api string, duration time.Duration) {
	_metricAPIRequestCount.With(group, p.Group, service, p.Service, serviceLabelType, api).Inc()
	_metricAPIRequestTime.With(group, p.Group, service, p.Service, serviceLabelType, api).Observe(float64(duration / time.Millisecond))
}

func (p ProcessMonitor) work() {
	load, _ := vload.Load(vload.OneMin)
	cpuUsage, _ := vcpu.Usage()

	_metricCPUUsage.With(group, p.Group, service, p.Service, instance, p.Instance).Set(cpuUsage)
	_metricLoadAvg1min.With(group, p.Group, service, p.Service, instance, p.Instance).Set(load)
	_metricMemory.With(group, p.Group, service, p.Service, instance, p.Instance).Set(float64(vmemory.Virtual()))
	_metricGoroutine.With(group, p.Group, service, p.Service, instance, p.Instance).Set(float64(vmemory.GoroutineNums()))
	_metricHeapObjects.With(group, p.Group, service, p.Service, instance, p.Instance).Set(float64(vmemory.HeapObjects()))
	_metricLastGCPause.With(group, p.Group, service, p.Service, instance, p.Instance).Set(vmemory.GCPause())
	_metricHeapAlloc.With(group, p.Group, service, p.Service, instance, p.Instance).Set(float64(vmemory.HeapAlloc()))
}
