package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vtool/vlog"
	"vtool/vprometheus/metric"
	"vtool/vprometheus/vcollector"
	"vtool/vservice/common"
	"vtool/vservice/server/engine"
	"vtool/vservice/server/register"
	"vtool/vservice/server/register/consul"
)

type ServiceBase struct {
	register       common.Register
	metricRegister common.Register

	baseLoc string
	name    string
	group   string
	ID      string

	servAddr string

	path string
	val  map[common.ServiceType]*common.ServiceInfo
	ttl  time.Duration

	shutDown func()
}

func NewServiceBase(ctx context.Context, args *servArgs) (*ServiceBase, error) {
	regEngine, err := register.GetRegisterEngine(args.registerType)
	if err != nil {
		return nil, err
	}

	return &ServiceBase{
		register:       regEngine,
		metricRegister: consul.DefaultConsulInstance,
		baseLoc:        common.DefaultRegisterPath,
		name:           args.serviceName,
		group:          args.serviceGroup,
		path:           common.DefaultRegisterPath + common.Slash + args.serviceGroup + common.Slash + args.serviceName,
		ttl:            common.DefaultTTl,
		shutDown: func() {
			vlog.InfoF(ctx, "service quit ~")
		},
	}, nil
}

func (sb *ServiceBase) Register(ctx context.Context, props map[common.ServiceType]common.Processor) error {
	serv, err := sb.powerServices(ctx, props)
	if err != nil {
		return err
	}

	val, err := json.Marshal(serv)
	if err != nil {
		return err
	}

	servID, err := sb.register.Register(ctx, sb.path, string(val), sb.ttl)
	if err != nil {
		return err
	}
	sb.ID = servID
	sb.val = serv
	sb.servAddr = string(val)

	sb.initMetric(ctx)
	return nil
}

func (sb *ServiceBase) initMetric(ctx context.Context) error {
	metric.InitBaseMetric(ctx, sb.group, sb.name, sb.servAddr)

	serv, err := sb.powerServices(ctx, map[common.ServiceType]common.Processor{
		common.Metric: &vcollector.MetricProcessor{},
	})

	if err != nil {
		return err
	}
	serviceInfo, ok := serv[common.Metric]
	if !ok {
		return nil
	}

	_, err = sb.metricRegister.Register(ctx, sb.FullServiceRegisterPath(), serviceInfo.Addr, common.DefaultTTl)
	if err != nil {
		return err
	}

	return nil
}

func (sb *ServiceBase) powerServices(ctx context.Context, props map[common.ServiceType]common.Processor) (map[common.ServiceType]*common.ServiceInfo, error) {
	serv := make(map[common.ServiceType]*common.ServiceInfo, len(props))

	for name, processor := range props {
		addr, engineFunc := processor.Engine()
		enginePower, ok := engine.GetEnginePower(engineFunc)
		if !ok {
			return nil, fmt.Errorf("not found engine power")
		}

		listenAddr, err := enginePower.Power(ctx, addr)
		if err != nil {
			return nil, err
		}

		serv[name] = &common.ServiceInfo{
			Type: enginePower.Type(),
			Addr: listenAddr,
		}
	}

	return serv, nil
}

func (sb *ServiceBase) ServName() string {
	return sb.name
}

func (sb *ServiceBase) ServGroup() string {
	return sb.group
}

func (sb *ServiceBase) ServInfo() map[common.ServiceType]*common.ServiceInfo {
	m := make(map[common.ServiceType]*common.ServiceInfo, len(sb.val))
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
	return sb.path + common.Slash + sb.ID
}

func (sb *ServiceBase) Stop() {
	ctx := context.Background()
	sb.register.UnRegister(ctx, sb.FullServiceRegisterPath())
	sb.metricRegister.UnRegister(ctx, sb.FullServiceRegisterPath())
	sb.shutDown()
}
