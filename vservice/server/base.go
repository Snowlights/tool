package server

import (
	"context"
	"vtool/vservice/common"
)

type ServiceBase struct {
	register       common.Register
	metricRegister common.Register

	baseLoc string
	name    string
	group   string
	ID      int64

	path string
	val  map[string]*common.ServiceInfo

	shutDown func()
}

func NewServiceBase() *ServiceBase {
	return &ServiceBase{}
}

func (sb *ServiceBase) ServName() string {
	return sb.name
}

func (sb *ServiceBase) ServGroup() string {
	return sb.group
}

func (sb *ServiceBase) ServInfo() map[string]*common.ServiceInfo {
	m := make(map[string]*common.ServiceInfo, len(sb.val))
	for k, v := range sb.val {
		m[k] = func() *common.ServiceInfo {
			return &common.ServiceInfo{
				Type: v.Type,
				Addr: v.Addr,
			}
		}()
	}
	return m
}

func (sb *ServiceBase) FullServiceRegisterPath() string {
	return sb.path
}

func (sb *ServiceBase) Stop() {
	ctx := context.Background()
	sb.register.UnRegister(ctx, sb.path)
	sb.metricRegister.UnRegister(ctx, sb.path)
}
