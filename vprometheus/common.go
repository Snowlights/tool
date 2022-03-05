package vprometheus

type VecOpts struct {
	NameSpace  string
	SubSystem  string
	Name       string
	Help       string
	LabelNames []string
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
