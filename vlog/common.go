package vlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FormatType int64

const (
	JsonFormatType FormatType = 1
	TextFormatType FormatType = 2
)

const (
	defaultLogMaxSize = 1024
)

var baseLogger *Logger

func init() {
	baseLogger = InitLogger("", "", zapcore.DebugLevel, JsonFormatType)
}

type Logger struct {
	zap   *zap.Logger
	level zapcore.Level

	ServiceName string
	Host        string
}
