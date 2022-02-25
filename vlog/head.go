package vlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	Body = "body"
	Msg  = "msg"
)

func (l Logger) buildHead(level, msg string) []zap.Field {
	kvs := append([]zap.Field{},
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
