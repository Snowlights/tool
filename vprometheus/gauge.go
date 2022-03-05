package vprometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"vtool/vlog"
)

type gaugeVec struct {
	gv  *prometheus.GaugeVec
	lvs LabelValues
}

func NewGauge(config *VecOpts) Gauge {
	gv := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: config.NameSpace,
		Subsystem: config.SubSystem,
		Name:      config.Name,
		Help:      config.Help,
	}, config.LabelNames)

	prometheus.MustRegister(gv)
	return &gaugeVec{gv: gv}
}

func (g *gaugeVec) Add(data float64) {
	if err := g.lvs.Check(); err != nil {
		vlog.Error(context.Background(), "gauge label value invalid:%s\n", err.Error())
		return
	}
	g.gv.With(makePrometheusLabels(g.lvs...)).Add(data)
}

func (g *gaugeVec) Sub(data float64) {
	if err := g.lvs.Check(); err != nil {
		vlog.Error(context.Background(), "gauge label value invalid:%s\n", err.Error())
		return
	}
	g.gv.With(makePrometheusLabels(g.lvs...)).Sub(data)
}

func (g *gaugeVec) Set(data float64) {
	if err := g.lvs.Check(); err != nil {
		vlog.Error(context.Background(), "gauge label value invalid:%s\n", err.Error())
		return
	}
	g.gv.With(makePrometheusLabels(g.lvs...)).Set(data)
}

func (g *gaugeVec) Inc() {
	if err := g.lvs.Check(); err != nil {
		vlog.Error(context.Background(), "gauge label value invalid:%s\n", err.Error())
		return
	}
	g.gv.With(makePrometheusLabels(g.lvs...)).Inc()
}

func (g *gaugeVec) Dec() {
	if err := g.lvs.Check(); err != nil {
		vlog.Error(context.Background(), "gauge label value invalid:%s\n", err.Error())
		return
	}
	g.gv.With(makePrometheusLabels(g.lvs...)).Dec()
}

func (g *gaugeVec) With(labelValues ...string) Gauge {
	return &gaugeVec{
		gv:  g.gv,
		lvs: g.lvs.With(labelValues...),
	}
}
