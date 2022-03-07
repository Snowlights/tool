package vprometheus

type VecOpts struct {
	NameSpace  string
	SubSystem  string
	Name       string
	Help       string
	LabelNames []string
	Buckets    []float64
}

type Counter interface {
	Inc()
	Add(float64)
	With(...string) Counter
}

type Gauge interface {
	Add(float64)
	Sub(float64)
	Set(float64)
	Inc()
	Dec()
	With(...string) Gauge
}

type Histogram interface {
	Observe(float64)
	With(...string) Histogram
}
