package vlog

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"time"
)

type (
	LogLevel zapcore.Level
)

const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	ErrorLevel = "ERROR"
	PanicLevel = "PANIC"
	FatalLevel = "FATAL"
)

const (
	Level = "level"

	Head            = "head"
	TimeHead        = "ts"
	ServiceNameHead = "service"
	ServiceHost     = "host"
	traceId         = "traceId"
	tracingError    = "error"

	callerFunction = "caller"

	Body = "body"
	Msg  = "msg"

	functionName = "?"
	character    = ":"
)

func (l Logger) getRuntimeInfo() (function, filename string, lineno int) {
	function = functionName
	pc, filename, lineno, ok := runtime.Caller(l.skip)
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}
	return
}

func (l Logger) buildHead(ctx context.Context, level, msg string) []zap.Field {

	// trace id
	var traceID string
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			traceID = fmt.Sprint(sc.TraceID())
		}

		// set tracing error log tag
		if level >= ErrorLevel {
			span.SetTag(tracingError, true)
		}
	}

	_, filename, lineno := l.getRuntimeInfo()
	if span != nil {
		span.LogKV(Level, level, callerFunction, fmt.Sprintf("%s:%d", filename, lineno),
			Body, msg)
	}
	kvs := append([]zap.Field{},
		zap.Field{
			Key:    callerFunction,
			Type:   zapcore.StringType,
			String: fmt.Sprintf("%s%s%d", filename, character, lineno),
		},
		zap.Field{
			Key:    Level,
			Type:   zapcore.StringType,
			String: level,
		},
		zap.Namespace(Head),
		zap.String(traceId, traceID),
		zap.Time(TimeHead, time.Now()),
		zap.String(ServiceNameHead, l.ServiceName),
		zap.String(ServiceHost, l.Host),
		zap.Namespace(Body),
		zap.String(Msg, msg),
	)

	return kvs
}
