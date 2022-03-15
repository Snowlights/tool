package vmetric

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestNewGauge(t *testing.T) {
	s := httptest.NewServer(promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	defer s.Close()

	scrape := func() string {
		resp, _ := http.Get(s.URL)
		buf, _ := ioutil.ReadAll(resp.Body)
		return string(buf)
	}

	namespace, subsystem, name := "namespace", "subsystem", "name"
	re := regexp.MustCompile(namespace + `_` + subsystem + `_` + name + `{a="a-value"} ([0-9\.]+)`)

	gauge := NewGauge(&VecOpts{
		NameSpace:  namespace,
		SubSystem:  subsystem,
		Name:       name,
		Help:       "help...",
		LabelNames: []string{"a"},
	}).With("a", "a-value")

	value := func() float64 {
		matches := re.FindStringSubmatch(scrape())
		f, _ := strconv.ParseFloat(matches[1], 64)
		return f
	}

	if err := gaugeTest(gauge, value); err != nil {
		t.Fatal(err)
	}

}

func gaugeTest(gauge Gauge, value func() float64) error {
	rand.Seed(time.Now().Unix())
	a := rand.Perm(100)
	n := rand.Intn(len(a))

	var want float64
	for i := 0; i < n; i++ {
		f := float64(a[i])
		gauge.Set(f)
		want = f
	}

	for i := 0; i < n; i++ {
		f := float64(a[i])
		gauge.Add(f)
		want += f
	}

	if have := value(); want != have {
		return fmt.Errorf("want %f, have %f", want, have)
	}

	return nil
}
