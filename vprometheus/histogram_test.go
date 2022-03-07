package vprometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHistogram(t *testing.T) {

	s := httptest.NewServer(promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	defer s.Close()

	scrape := func() string {
		resp, _ := http.Get(s.URL)
		buf, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(buf))
		return string(buf)
	}

	namespace, subsystem, name := "namespace", "subsystem", "name"
	// re := regexp.MustCompile(namespace + `_` + subsystem + `_` + name + `_bucket{x="1",le="([0-9]+|\+Inf)"} ([0-9\.]+)`)

	bucketMin, bucketMax, bucketStep := 0, 100, 10
	buckets := []float64{}
	for i := bucketMin; i <= bucketMax; i += bucketStep {
		buckets = append(buckets, float64(i))
	}

	histogram := NewHistogram(&VecOpts{
		NameSpace:  namespace,
		SubSystem:  subsystem,
		Name:       name,
		Help:       "help",
		LabelNames: []string{"a"},
		Buckets:    buckets,
	}).With("a", "1")

	SetHistogram(histogram)
	SetHistogram(histogram)

	fmt.Println(scrape())
}

func SetHistogram(h Histogram) {
	for i := 0; i < 100; i++ {
		h.Observe(float64(i))
	}
}
