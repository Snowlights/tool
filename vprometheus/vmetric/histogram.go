package vmetric

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/Snowlights/tool/vlog"
)

// avg(http_request_duration_seconds{quantile="0.95"}) // BAD!
// histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) // GOOD.

type histogramVec struct {
	hv  *prometheus.HistogramVec
	lvs LabelValues
}

func NewHistogram(config *VecOpts) Histogram {
	hv := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.NameSpace,
		Subsystem: config.SubSystem,
		Name:      config.Name,
		Help:      config.Help,
		Buckets:   config.Buckets,
	}, config.LabelNames)

	prometheus.MustRegister(hv)
	return &histogramVec{hv: hv}
}

func (hv *histogramVec) Observe(data float64) {
	if err := hv.lvs.Check(); err != nil {
		vlog.ErrorF(context.Background(), "histogram label value invalid:%s", err.Error())
		return
	}

	hv.hv.With(makePrometheusLabels(hv.lvs...)).Observe(data)
}

func (hv *histogramVec) With(labelValues ...string) Histogram {
	return &histogramVec{
		hv:  hv.hv,
		lvs: hv.lvs.With(labelValues...),
	}
}
