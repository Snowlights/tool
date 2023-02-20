package vmetric

import (
	"context"
	"github.com/Snowlights/tool/vlog"
	"github.com/prometheus/client_golang/prometheus"
)

type counterVec struct {
	cv  *prometheus.CounterVec
	lvs LabelValues
}

func NewCounter(config *VecOpts) Counter {
	cv := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.NameSpace,
		Subsystem: config.SubSystem,
		Name:      config.Name,
		Help:      config.Help,
	}, config.LabelNames)

	prometheus.MustRegister(cv)
	return &counterVec{cv: cv}
}

// Inc Inc increments the counter by 1.
func (c *counterVec) Inc() {
	if err := c.lvs.Check(); err != nil {
		vlog.ErrorF(context.Background(), "counter label value invalid:%s", err.Error())
		return
	}
	c.cv.With(makePrometheusLabels(c.lvs...)).Inc()
}

// Add adds the given value to the counter. It panics if the value is < 0
func (c *counterVec) Add(data float64) {
	if err := c.lvs.Check(); err != nil {
		vlog.ErrorF(context.Background(), "counter label value invalid:%s", err.Error())
		return
	}
	c.cv.With(makePrometheusLabels(c.lvs...)).Add(data)
}

// With implements Counter, labels will see as a key
func (c *counterVec) With(labelValues ...string) Counter {
	return &counterVec{
		cv:  c.cv,
		lvs: c.lvs.With(labelValues...),
	}
}
