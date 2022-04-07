package common

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"vtool/idl/thrift/gen-go/thriftBase"
	"vtool/vservice/common"
	"vtool/vservice/server"
)

const (
	InjectServKey = "InjectServInfo"
)

func NewContextFromThriftBaseContext(operation string, tctx *thriftBase.Context) context.Context {
	ctx := context.Background()

	tracer := opentracing.GlobalTracer()
	spanCtx, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(tctx.SpanCtx))
	var span opentracing.Span
	if err == nil {
		span = tracer.StartSpan(operation, ext.RPCServerOption(spanCtx))
	} else {
		span = tracer.StartSpan(operation)
	}

	if tctx.SpanCtx != nil {
		ctx = context.WithValue(ctx, InjectServKey, tctx.SpanCtx)
	}

	servBase := server.GetServBase()
	if servBase != nil {
		servInfo := servBase.ServInfo()
		if servInfo != nil {
			span.SetTag("lane", servInfo.Lane)
			serv, ok := servInfo.ServList[common.Grpc]
			if ok {
				span.SetTag("servType", common.Grpc)
				span.SetTag("servIp", serv.Addr)
				span.SetTag("engineType", serv.Type)
			}
			serv, ok = servInfo.ServList[common.Thrift]
			if ok {
				span.SetTag("component", "thrift")
				span.SetTag("span.kind", "server")
				span.SetTag("servType", common.Thrift)
				span.SetTag("servIp", serv.Addr)
				span.SetTag("engineType", serv.Type)
			}
		}
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx
}

func NewThriftBaseContextFromContext(ctx context.Context) *thriftBase.Context {

	if ctx == nil {
		ctx = context.Background()
	}

	tctx := &thriftBase.Context{
		SpanCtx: make(map[string]string),
	}

	carrier := opentracing.TextMapCarrier(make(map[string]string))
	span := opentracing.SpanFromContext(ctx)

	if span != nil {
		opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.TextMap,
			carrier)
	}
	tctx.SpanCtx = carrier

	return tctx
}
