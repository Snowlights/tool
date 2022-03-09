package vcollector

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"vtool/server/consul"
)

//Nerve SD configurations allow retrieving scrape targets from [AirBnB's Nerve]
// (https://github.com/airbnb/nerve) which are stored in Zookeeper.

// Serverset SD configurations allow retrieving scrape targets from
// [Serversets] (https://github.com/twitter/finagle/tree/master/finagle-serversets)
// which are stored in Zookeeper. Serversets are commonly used by Finagle and Aurora.

// Zookeeper only supports these two kinds of structured data,
// but we implement service registration and discovery by ourselves, so we don't use ZK for service index statistics

const DefaultMetricPath = "/metric"

type MetricProcessor struct{}

func (mp *MetricProcessor) Prepare() error {
	// set default metric

	return nil
}

func (mp *MetricProcessor) Engine() (string, interface{}) {

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET(DefaultMetricPath, func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	engine.GET(consul.DefaultCheckPath, func(c *gin.Context) {
		consul.Checker{}.ServeHTTP(c.Writer, c.Request)
	})

	return "", engine
}
