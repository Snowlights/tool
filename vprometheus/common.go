package vprometheus

type CounterVec struct {
	NameSpace  string
	SubSystem  string
	Name       string
	Help       string
	LabelNames []string
}

type Counter interface {
	Inc()
	Add(data float64)
	With(labelValues ...string) Counter
}
