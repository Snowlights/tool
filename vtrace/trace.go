package vtrace

import (
	"github.com/Snowlights/tool/vconfig"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"os"
	"strconv"
	"sync"
)

// https://opentracing.io/guides/golang/quick-start/

var GlobalTracer *JaegerTracer

type JaegerTracer struct {
	servName string

	center vconfig.Center

	mu         sync.Mutex
	baseTracer *baseTracer
}

type baseTracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func (bt *baseTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return bt.tracer.StartSpan(operationName, opts...)
}
func (bt *baseTracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return bt.tracer.Inject(sm, format, carrier)
}

func (bt *baseTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return bt.tracer.Extract(format, carrier)
}

func InitJaegerTracer(servName string) error {
	jt := &JaegerTracer{
		servName: servName,
	}

	err := jt.initCenter()
	if err != nil {
		return err
	}

	cfg := jt.buildJaegerConfig()
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return err
	}

	jt.updateTracer(tracer, closer)
	// set global tracer
	opentracing.SetGlobalTracer(jt.baseTracer)
	GlobalTracer = jt
	return nil
}

func (jt *JaegerTracer) GetTracer() *baseTracer {
	return jt.baseTracer
}

func (jt *JaegerTracer) updateTracer(tracer opentracing.Tracer, closer io.Closer) {
	if jt.baseTracer == nil {
		jt.baseTracer = &baseTracer{}
	}

	jt.mu.Lock()
	defer jt.mu.Unlock()

	if jt.baseTracer.closer != nil {
		jt.baseTracer.closer.Close()
	}
	jt.baseTracer.tracer = tracer
	jt.baseTracer.closer = closer
}

func (jt *JaegerTracer) Close() {
	jt.baseTracer.closer.Close()
}

func (jt *JaegerTracer) buildJaegerConfig() *config.Configuration {
	return &config.Configuration{
		ServiceName: jt.servName,
		Sampler:     jt.buildSampler(),
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: jt.getLocalEngineHost(),
		},
		Headers: jt.builderHeaders(),
	}
}

func (jt *JaegerTracer) buildSampler() (res *config.SamplerConfig) {
	// think about namespaces and services need to be sampled
	res = defaultSampler

	servSampler := jt.getServSampler()
	if servSampler != nil {
		res = servSampler
		return
	}

	globalSampler := jt.getGlobalSampler()
	if globalSampler != nil {
		res = globalSampler
		return
	}

	return
}

func (jt *JaegerTracer) getServSampler() *config.SamplerConfig {
	typeKey, paramKey := buildServSamplerTypeKey(jt.servName), buildServSamplerParamKey(jt.servName)
	typeKeyVal, ok := jt.center.GetValueWithNamespace(vconfig.MiddlewareNamespaceTrace, typeKey)
	if !ok {
		return nil
	}
	paramKeyVal, ok := jt.center.GetValueWithNamespace(vconfig.MiddlewareNamespaceTrace, paramKey)
	if !ok {
		return nil
	}
	samplerParam, err := strconv.ParseFloat(paramKeyVal, 64)
	if err != nil {
		return nil
	}

	return &config.SamplerConfig{
		Type:  typeKeyVal,
		Param: samplerParam,
	}
}

func (jt *JaegerTracer) getGlobalSampler() *config.SamplerConfig {

	typeKey, paramKey := globalSamplerType, globalSamplerParam
	typeKeyVal, ok := jt.center.GetValueWithNamespace(vconfig.MiddlewareNamespaceTrace, typeKey)
	if !ok {
		return nil
	}
	paramKeyVal, ok := jt.center.GetValueWithNamespace(vconfig.MiddlewareNamespaceTrace, paramKey)
	if !ok {
		return nil
	}
	samplerParam, err := strconv.ParseFloat(paramKeyVal, 64)
	if err != nil {
		return nil
	}
	return &config.SamplerConfig{
		Type:  typeKeyVal,
		Param: samplerParam,
	}
}

func (jt *JaegerTracer) getLocalEngineHost() string {

	agentHost, agentPort := defaultAgentHost, defaultAgentPort

	if host, ok := os.LookupEnv(JaegerAgentHost); ok {
		agentHost = host
	}

	if port, ok := os.LookupEnv(JaegerAgentPort); ok {
		agentPort = port
	}

	return agentHost + colon + agentPort
}

func (jt *JaegerTracer) builderHeaders() *jaeger.HeadersConfig {
	return &jaeger.HeadersConfig{
		JaegerDebugHeader:        TraceDebugHeader,
		JaegerBaggageHeader:      TraceBaggageHeader,
		TraceContextHeaderName:   TraceContextHeaderName,
		TraceBaggageHeaderPrefix: TraceBaggageHeaderPrefix,
	}
}

func (jt *JaegerTracer) initCenter() error {
	cfg, err := vconfig.ParseConfigEnv()
	if err != nil {
		panic(err)
	}

	port, err := strconv.ParseInt(cfg.Port, 10, 64)
	if err != nil {
		panic(err)
	}

	center, err := vconfig.NewCenter(&vconfig.CenterConfig{
		AppID:     vconfig.MiddlewareAppID,
		Cluster:   cfg.Cluster,
		Namespace: []string{vconfig.MiddlewareNamespaceTrace},
		IP:        cfg.IP,
		Port:      int(port),
		MustStart: true,
	})
	if err != nil {
		return err
	}

	jt.center = center
	return nil
}
