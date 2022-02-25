package vlog

import (
	"context"
	"go.uber.org/zap/zapcore"
)

func NewLogger(rootPath, fileName string, logLevel zapcore.Level, formatType FormatType) *Logger {
	return InitLogger(rootPath, fileName, logLevel, formatType)
}

func Debug(ctx context.Context, args ...interface{}) {
	baseLogger.debug(ctx, args...)
}

func DebugF(ctx context.Context, template string, args ...interface{}) {
	baseLogger.debugF(ctx, template, args...)
}

func Info(ctx context.Context, args ...interface{}) {

}

func InfoF(ctx context.Context, template string, args ...interface{}) {

}

func Warn(ctx context.Context, args ...interface{}) {

}

func WarnF(ctx context.Context, template string, args ...interface{}) {

}

func Error(ctx context.Context, args ...interface{}) {

}

func ErrorF(ctx context.Context, template string, args ...interface{}) {

}

func Fatal(ctx context.Context, args ...interface{}) {

}

func FatalF(ctx context.Context, template string, args ...interface{}) {

}

func Panic(ctx context.Context, args ...interface{}) {

}

func PanicF(ctx context.Context, template string, args ...interface{}) {

}
