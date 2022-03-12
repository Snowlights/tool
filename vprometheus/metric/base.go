package metric

import "vtool/vprometheus/vmetric"

const (
	namespace = "process_runtime_resource"
	group     = "group"
	service   = "service"
	instance  = "instance"
)

var labelList = []string{group, service, instance}

var (
	// cpu usage
	_metricCPUUsage = vmetric.NewGauge(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  "cpu_usage",
		Name:       "current",
		Help:       "cup usage percentage",
		LabelNames: labelList,
	})
	// load avg 1 min
	_metricLoadAvg1min = vmetric.NewGauge(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  "load_avg",
		Name:       "current",
		Help:       "load avg 1 min",
		LabelNames: labelList,
	})
	// memory of process
	_metricMemory = vmetric.NewGauge(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  "memory",
		Name:       "current",
		Help:       "memory of process, data from runtime memory.sys",
		LabelNames: labelList,
	})
	// goroutine number
	_metricGoroutine = vmetric.NewGauge(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  "goroutine",
		Name:       "current",
		Help:       "runtime numGoroutine",
		LabelNames: labelList,
	})
	// heap objects number
	_metricHeapObjects = vmetric.NewGauge(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  "heap_objects",
		Name:       "current",
		Help:       "runtime heapObjects",
		LabelNames: labelList,
	})
	// last gc pause
	_metricLastGCPause = vmetric.NewGauge(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  "gc_pause",
		Name:       "last",
		Help:       "last gc pause time",
		LabelNames: labelList,
	})
	// heap alloc
	_metricHeapAlloc = vmetric.NewGauge(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  "heap_alloc",
		Name:       "current",
		Help:       "runtime heapAlloc",
		LabelNames: labelList,
	})
)
