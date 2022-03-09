package vcollector

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
	"time"
	"vtool/vnet"
	"vtool/vprometheus/vmetric"
)

func TestNewCollector(t *testing.T) {

	namespace, subsystem, name := "namespace", "subsystem", "name"

	counter := vmetric.NewCounter(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  subsystem,
		Name:       name,
		Help:       "help info",
		LabelNames: []string{"a", "b"},
	}).With("a", "a-value", "b", "b-value")

	go func() {
		for {
			counter.Add(1)
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}

func TestNewCollector2(t *testing.T) {

	namespace, subsystem, name := "namespace", "subsystem", "gauge"

	gauge := vmetric.NewGauge(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  subsystem,
		Name:       name,
		Help:       "help...",
		LabelNames: []string{"a"},
	}).With("a", "a-value")

	go func() {
		for {
			gauge.Add(1)
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func TestNewCollector3(t *testing.T) {

	namespace, subsystem, name := "namespace", "subsystem", "histogram"

	bucketMin, bucketMax, bucketStep := 0, 100, 10
	buckets := []float64{}
	for i := bucketMin; i <= bucketMax; i += bucketStep {
		buckets = append(buckets, float64(i))
	}

	histogram := vmetric.NewHistogram(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  subsystem,
		Name:       name,
		Help:       "help",
		LabelNames: []string{"a"},
		Buckets:    buckets,
	}).With("a", "1")

	go func() {
		for {
			for i := 0; i < 100; i++ {
				histogram.Observe(float64(i))
			}
			time.Sleep(time.Second * 15)
		}
	}()

	listen, err := vnet.ListenServAddr(context.Background(), ":")
	if err != nil {
		return
	}

	http.Handle("/health", MyHandler{})
	http.Handle("/metrics", promhttp.Handler())
	http.Serve(listen, nil)
}

type MyHandler struct{}

func (MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("heath check success!"))
}
