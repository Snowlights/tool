package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"vtool/vlog"
	"vtool/vprometheus/vcollector"
	"vtool/vservice/common"
	"vtool/vservice/service/engine"
	register2 "vtool/vservice/service/register"
	"vtool/vservice/service/register/consul"
)

func Serv(ctx context.Context, registerConfig *common.RegisterConfig, props map[string]common.Processor) error {
	err := serverIns(ctx, registerConfig, props)
	if err != nil {
		return err
	}
	awaitSignal()
	return nil
}

func awaitSignal() {
	c := make(chan os.Signal, 1)
	ctx := context.Background()
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGPIPE}
	signal.Reset(signals...)
	signal.Notify(c, signals...)

	for {
		select {
		case s := <-c:
			vlog.InfoF(ctx, "receive a signal:%s", s.String())

			if s.String() == syscall.SIGTERM.String() {
				vlog.InfoF(ctx, "receive a signal: %s, stop vservice", s.String())
				Stop()
				<-(chan int)(nil)
			}
		}
	}
}

func Stop() {
	// clearRegisterInfos
}

func serverIns(ctx context.Context, registerConfig *common.RegisterConfig, props map[string]common.Processor) error {

	// power service
	serv, err := powerServices(ctx, props)
	if err != nil {
		return err
	}
	// register service
	err = register(ctx, registerConfig, serv)
	if err != nil {
		return err
	}

	initMetric(ctx, registerConfig)

	return nil
}

func register(ctx context.Context, registerConfig *common.RegisterConfig, serv map[string]*common.ServiceInfo) error {
	servStr, err := json.Marshal(serv)
	if err != nil {
		return err
	}

	regConfig := &common.RegisterConfig{
		RegistrationType: registerConfig.RegistrationType,
		ServName:         registerConfig.ServName,
		ServAddr:         string(servStr),
		Group:            registerConfig.Group,
	}

	err = register2.RegisterService(ctx, regConfig)
	if err != nil {
		return err
	}

	return nil
}

func initMetric(ctx context.Context, registerConfig *common.RegisterConfig) error {

	serv, err := powerServices(ctx, map[string]common.Processor{
		"metric": &vcollector.MetricProcessor{},
	})
	if err != nil {
		return err
	}
	serviceInfo, ok := serv["metric"]
	if !ok {
		return nil
	}

	err = consul.DefaultConsulInstance.Register(ctx, consul.ConsulNamespace+strings.Join([]string{registerConfig.Group, registerConfig.ServName, "0"}, common.Slash),
		serviceInfo.Addr, common.DefaultTTl)
	if err != nil {
		return err
	}

	return nil
}

func powerServices(ctx context.Context, props map[string]common.Processor) (map[string]*common.ServiceInfo, error) {
	serv := make(map[string]*common.ServiceInfo, len(props))

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
