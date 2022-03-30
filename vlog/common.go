package vlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FormatType int64

const (
	JsonFormatType FormatType = 1
	TextFormatType FormatType = 2

	LogPath = "logs"
	LogFile = "vlog.log"
)

const (
	defaultLogMaxSize = 1024
	defaultSkip       = 4
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

	skip int
}
