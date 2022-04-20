package vtrace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"net/http"
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

const (
	Component = "component"
	SpanKind  = "span.kind"
)

// server span tags
const (
	Lane       = "lane"
	ServType   = "servType"
	ServIP     = "servIP"
	EngineType = "engineType"
)

// sql span tags
const (
	ComponentSQL = "sql"
	SpanKindSQL  = "client"

	Cluster = "cluster"
	Schema  = "schema"
	Table   = "table"
	Query   = "query"
)

// redis span tags
const (
	ComponentRedis = "redis"
	SpanKindRedis  = "client"

	RedisCluster = "cluster"
	RedisCmd     = "cmd"
	WithPipeline = "withPipeline"
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

func TraceHTTPRequest(ctx context.Context, req *http.Request) error {
	if ctx == nil {
		return fmt.Errorf("TraceHTTPRequest got nil context")
	}

	if req == nil {
		return fmt.Errorf("TraceHTTPRequest got nil request")
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		return opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
	}

	return fmt.Errorf("TraceHTTPRequest got nil span")
}
