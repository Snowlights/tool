package vcollector

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
	"github.com/Snowlights/tool/vservice/server/register/consul"
)

//Nerve SD configurations allow retrieving scrape targets from [AirBnB's Nerve]
// (https://github.com/airbnb/nerve) which are stored in Zookeeper.

// Serverset SD configurations allow retrieving scrape targets from
// [Serversets] (https://github.com/twitter/finagle/tree/master/finagle-serversets)
// which are stored in Zookeeper. Serversets are commonly used by Finagle and Aurora.

// Zookeeper only supports these two kinds of structured data,
// but we implement server registration and discovery by ourselves, so we don't use ZK for server index statistics

//var profileDescriptions = map[string]string{
//	"allocs":       "A sampling of all past memory allocations",
//	"block":        "Stack traces that led to blocking on synchronization primitives",
//	"cmdline":      "The command line invocation of the current program",
//	"goroutine":    "Stack traces of all current goroutines",
//	"heap":         "A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.",
//	"mutex":        "Stack traces of holders of contended mutexes",
//	"profile":      "CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.",
//	"threadcreate": "Stack traces that led to the creation of new OS threads",
//	"trace":        "A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.",
//}
const (
	DefaultMetricPath = "/metrics"
	pprofPath         = "/debug/pprof/"
	allocs            = "allocs"
	block             = "block"
	cmdline           = "cmdline"
	goroutine         = "goroutine"
	heap              = "heap"
	mutex             = "mutex"
	threadcreate      = "threadcreate"
	trace             = "trace"
	profile           = "profile"
	symbol            = "symbol"
)

type MetricProcessor struct{}

func (mp *MetricProcessor) Engine() (string, interface{}) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET(DefaultMetricPath, func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})
	engine.GET(consul.DefaultCheckPath, func(c *gin.Context) {
		consul.Checker{}.ServeHTTP(c.Writer, c.Request)
	})

	// pprof handler, use for trace info later
	engine.GET(pprofPath, pprofHandler(pprof.Index))
	engine.POST(pprofPath, pprofHandler(pprof.Index))
	engine.GET(pprofPath+allocs, pprofHandler(pprof.Index))
	engine.POST(pprofPath+allocs, pprofHandler(pprof.Index))
	engine.GET(pprofPath+block, pprofHandler(pprof.Index))
	engine.POST(pprofPath+block, pprofHandler(pprof.Index))
	engine.GET(pprofPath+cmdline, pprofHandler(pprof.Cmdline))
	engine.POST(pprofPath+cmdline, pprofHandler(pprof.Cmdline))
	engine.GET(pprofPath+goroutine, pprofHandler(pprof.Index))
	engine.POST(pprofPath+goroutine, pprofHandler(pprof.Index))
	engine.GET(pprofPath+heap, pprofHandler(pprof.Index))
	engine.POST(pprofPath+heap, pprofHandler(pprof.Index))
	engine.GET(pprofPath+mutex, pprofHandler(pprof.Index))
	engine.POST(pprofPath+mutex, pprofHandler(pprof.Index))
	engine.GET(pprofPath+threadcreate, pprofHandler(pprof.Index))
	engine.POST(pprofPath+threadcreate, pprofHandler(pprof.Index))
	engine.GET(pprofPath+trace, pprofHandler(pprof.Trace))
	engine.POST(pprofPath+trace, pprofHandler(pprof.Trace))
	engine.GET(pprofPath+profile, pprofHandler(pprof.Profile))
	engine.POST(pprofPath+profile, pprofHandler(pprof.Profile))
	engine.GET(pprofPath+symbol, pprofHandler(pprof.Symbol))
	engine.POST(pprofPath+symbol, pprofHandler(pprof.Symbol))

	return "", engine
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
