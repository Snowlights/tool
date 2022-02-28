package vlog

import (
	"fmt"
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

func (l Logger) buildHead(level, msg string) []zap.Field {

	_, filename, lineno := l.getRuntimeInfo()
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
		zap.Time(TimeHead, time.Now()),
		zap.String(ServiceNameHead, l.ServiceName),
		zap.String(ServiceHost, l.Host),
		zap.Namespace(Body),
		zap.String(Msg, msg),
	)

	return kvs
}
