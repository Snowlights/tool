package vconfig

import (
	"context"
	"vtool/vlog"
)

type CenterLogger struct {
	c context.Context
}

func (cl *CenterLogger) Debugf(format string, params ...interface{}) {
	vlog.DebugF(cl.c, format, params...)
}

func (cl *CenterLogger) Infof(format string, params ...interface{}) {
	vlog.InfoF(cl.c, format, params...)
}

func (cl *CenterLogger) Warnf(format string, params ...interface{}) {
	vlog.WarnF(cl.c, format, params...)
}

func (cl *CenterLogger) Errorf(format string, params ...interface{}) {
	vlog.ErrorF(cl.c, format, params...)
}

func (cl *CenterLogger) Debug(v ...interface{}) {
	vlog.Debug(cl.c, v...)
}

func (cl *CenterLogger) Info(v ...interface{}) {
	vlog.Info(cl.c, v...)
}

func (cl *CenterLogger) Warn(v ...interface{}) {
	vlog.Warn(cl.c, v...)
}

func (cl *CenterLogger) Error(v ...interface{}) {
	vlog.Error(cl.c, v...)
}
