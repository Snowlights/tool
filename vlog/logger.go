package vlog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func InitLogger(path string, name string, level zapcore.Level, formatType FormatType) *Logger {
	baseLogger = &Logger{
		level: level,
	}

	logFile := ""
	if path != "" && name != "" {
		logFile = path + "/" + name
	}
	baseLogger.level = level

	var w io.Writer
	if logFile != "" {
		ljWriter := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    defaultLogMaxSize,
			MaxAge:     0,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   false,
		}
		go func() {
			for {
				now := time.Now().Unix()
				duration := 3600 - now%3600
				select {
				case <-time.After(time.Second * time.Duration(duration)):
					ljWriter.Rotate()
				}
			}
		}()
		w = ljWriter
	} else {
		w = os.Stdout
	}

	encodeConfig := zapcore.EncoderConfig{
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var core zapcore.Core
	switch formatType {
	case JsonFormatType, TextFormatType:
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encodeConfig),
			zapcore.AddSync(w), baseLogger.level)
	}

	baseLogger.zap = zap.New(core)
	return baseLogger
}

func (l *Logger) checkIsLog(level zapcore.Level) bool {
	return l.level >= level
}

func (l Logger) debug(ctx context.Context, args ...interface{}) {
	if !l.checkIsLog(zapcore.DebugLevel) {
		return
	}
	l.zap.Debug(fmt.Sprint(args...), l.buildHead(DebugLevel, fmt.Sprint(args...))...)
}

func (l Logger) debugF(ctx context.Context, template string, args ...interface{}) {
	if !l.checkIsLog(zapcore.DebugLevel) {
		return
	}
	l.zap.Debug(fmt.Sprintf(template, args...))
}

func (l Logger) info(ctx context.Context, args ...interface{}) {
	if !l.checkIsLog(zapcore.DebugLevel) {
		return
	}
	l.zap.Info(fmt.Sprint(args...))
}

func (l Logger) infoF(ctx context.Context, template string, args ...interface{}) {
	if !l.checkIsLog(zapcore.DebugLevel) {
		return
	}
	l.zap.Info(fmt.Sprintf(template, args...))
}
