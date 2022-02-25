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
	baseLogger.info(ctx, args...)
}

func InfoF(ctx context.Context, template string, args ...interface{}) {
	baseLogger.infoF(ctx, template, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	baseLogger.warn(ctx, args...)
}

func WarnF(ctx context.Context, template string, args ...interface{}) {
	baseLogger.warnF(ctx, template, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	baseLogger.error(ctx, args...)
}

func ErrorF(ctx context.Context, template string, args ...interface{}) {
	baseLogger.errorF(ctx, template, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	baseLogger.fatal(ctx, args...)
}

func FatalF(ctx context.Context, template string, args ...interface{}) {
	baseLogger.fatalF(ctx, template, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	baseLogger.panic(ctx, args...)
}

func PanicF(ctx context.Context, template string, args ...interface{}) {
	baseLogger.panicF(ctx, template, args...)
}
