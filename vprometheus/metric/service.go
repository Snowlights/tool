package metric

import (
	"time"
	"vtool/vprometheus/vmetric"
)

const (
	serviceLabelType = "servLabelType"

	apiType         = "api"
	count           = "requestCount"
	requestDuration = "requestDuration"
)

var (
	buckets = []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}

	_metricAPIRequestCount = vmetric.NewCounter(&vmetric.VecOpts{
		NameSpace:  serviceNamespace,
		SubSystem:  apiType,
		Name:       count,
		Help:       "api request count",
		LabelNames: []string{group, service, serviceLabelType},
	})

	_metricAPIRequestTime = vmetric.NewHistogram(&vmetric.VecOpts{
		NameSpace:  serviceNamespace,
		SubSystem:  apiType,
		Name:       requestDuration,
		Buckets:    buckets,
		Help:       "api request time",
		LabelNames: []string{group, service, serviceLabelType},
	})
)

func StatApi(api string, duration time.Duration) {
	defaultProcessMonitor.statApi(api, duration)
}
