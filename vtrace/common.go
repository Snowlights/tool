package vtrace

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	defaultAgentHost = "127.0.0.1"
	defaultAgentPort = "6831"

	JaegerAgentHost = "JAEGER_AGENT_HOST"
	JaegerAgentPort = "JAEGER_AGENT_PORT"

	colon = ":"
)

const (
	TraceDebugHeader         = "trace-debug-id"
	TraceBaggageHeader       = "trace-baggage-header"
	TraceContextHeaderName   = "trace-context-id"
	TraceBaggageHeaderPrefix = "trace-baggage-header-prefix"
)

const (
	servSamplerTypePrefix  = "serv_sampler_type."
	servSamplerParamPrefix = "serv_sampler_param."

	globalSamplerType  = "global_sampler_type"
	globalSamplerParam = "global_sampler_param"
)

func buildServSamplerTypeKey(servName string) string {
	return servSamplerTypePrefix + servName
}

func buildServSamplerParamKey(servName string) string {
	return servSamplerParamPrefix + servName
}

var defaultSampler = &config.SamplerConfig{
	Type:  jaeger.SamplerTypeRateLimiting,
	Param: 1.0,
}

func SpanFromContent(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}
