package metric

import (
	"github.com/Snowlights/tool/vprometheus/vmetric"
	"time"
)

const (
	serviceLabelType = "servLabelType"

	apiType         = "api"
	requestCount    = "requestCount"
	requestDuration = "requestDuration"
)

var (
	msBuckets = []float64{1, 5, 10, 25, 50, 100, 200, 300, 500, 1000, 3000, 5000, 10000, 15000}

	_metricAPIRequestCount = vmetric.NewCounter(&vmetric.VecOpts{
		NameSpace:  serviceNamespace,
		SubSystem:  apiType,
		Name:       requestCount,
		Help:       "api request count",
		LabelNames: []string{group, service, serviceLabelType},
	})

	_metricAPIRequestTime = vmetric.NewHistogram(&vmetric.VecOpts{
		NameSpace:  serviceNamespace,
		SubSystem:  apiType,
		Name:       requestDuration,
		Buckets:    msBuckets,
		Help:       "api request time",
		LabelNames: []string{group, service, serviceLabelType},
	})
)

func StatApi(api string, duration time.Duration) {
	defaultProcessMonitor.statApi(api, duration)
}
