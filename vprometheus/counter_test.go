package vprometheus

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

func TestNewCounter(t *testing.T) {
	s := httptest.NewServer(promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	defer s.Close()

	scrape := func() string {
		resp, _ := http.Get(s.URL)
		buf, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(buf))
		return string(buf)
	}

	namespace, subsystem, name := "namespace", "subsystem", "name"
	re := regexp.MustCompile(namespace + `_` + subsystem + `_` + name + `{a="a-value",b="b-value"} ([0-9\.]+)`)

	counter := NewCounter(&VecOpts{
		NameSpace:  namespace,
		SubSystem:  subsystem,
		Name:       name,
		Help:       "help info",
		LabelNames: []string{"a", "b"},
	}).With("a", "a-value", "b", "b-value")

	value := func() float64 {
		matches := re.FindStringSubmatch(scrape())
		f, _ := strconv.ParseFloat(matches[1], 64)
		return f
	}

	if err := count(counter, value); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("success ")
	}
}

func FillCounter(counter Counter) float64 {
	rand.Seed(time.Now().Unix())
	a := rand.Perm(100)
	n := rand.Intn(len(a))

	var want float64
	for i := 0; i < n; i++ {
		f := float64(a[i])
		counter.Add(f)
		want += f
	}
	return want
}

func count(counter Counter, value func() float64) error {
	want := FillCounter(counter)
	if have := value(); want != have {
		return fmt.Errorf("want %f, have %f", want, have)
	}

	return nil
}
