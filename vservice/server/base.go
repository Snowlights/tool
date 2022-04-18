package server

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
	"vtool/parse"
	"vtool/vconfig"
	"vtool/vlog"
	"vtool/vprometheus/metric"
	"vtool/vprometheus/vcollector"
	"vtool/vservice/common"
	"vtool/vservice/server/engine"
	"vtool/vservice/server/register"
	"vtool/vservice/server/register/consul"
	"vtool/vsql"
	"vtool/vtrace"
)

type ServiceBase struct {
	center vconfig.Center

	// todo add redis open api
	// todo add mq open api

	// todo: log log time, add region and cross region and colony config

	register       common.Register
	metricRegister common.Register

	baseLoc string
	name    string
	group   string
	lane    string
	ID      string

	servAddr string

	md5     string
	version int64

	path string
	val  *common.RegisterServiceInfo
	ttl  time.Duration

	shutDown func()
}

func NewServiceBase(ctx context.Context, args *servArgs) (*ServiceBase, error) {

	servBase := &ServiceBase{
		baseLoc: common.DefaultRegisterPath,
		name:    args.serviceName,
		group:   args.serviceGroup,
		lane:    args.serviceLane,
		version: args.version,
		path:    common.DefaultRegisterPath + common.Slash + args.serviceGroup + common.Slash + args.serviceName,
		ttl:     common.DefaultTTl,
		shutDown: func() {
			vlog.InfoF(ctx, "service quit ~")
		},
	}

	err := servBase.initCenter(args)
	if err != nil {
		return nil, err
	}

	_, err = vsql.InitManager(servBase.center)
	if err != nil {
		return nil, err
	}

	err = servBase.initRegisterEngines()
	if err != nil {
		return nil, err
	}

	vtrace.InitJaegerTracer(servBase.ServName())

	return servBase, nil
}

func (sb *ServiceBase) GetCenter(ctx context.Context) vconfig.Center {
	return sb.center
}

func (sb *ServiceBase) Register(ctx context.Context, props map[common.ServiceType]common.Processor) error {
	serv, err := sb.powerServices(ctx, props)
	if err != nil {
		return err
	}

	serv.ServPath = sb.path
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

func (sb *ServiceBase) ServName() string {
	return sb.name
}

func (sb *ServiceBase) ServGroup() string {
	return sb.group
}

func (sb *ServiceBase) ServInfo() *common.RegisterServiceInfo {
	f := new(common.RegisterServiceInfo)
	f = sb.val
	return f
}

func (sb *ServiceBase) FullServiceRegisterPath() string {
	return sb.path + common.Slash + sb.ID
}

func (sb *ServiceBase) Stop() {
	ctx := context.Background()
	if sb.register != nil {
		sb.register.UnRegister(ctx, sb.FullServiceRegisterPath())
	}
	if sb.metricRegister != nil {
		sb.metricRegister.UnRegister(ctx, sb.FullServiceRegisterPath())
	}
	if sb.shutDown != nil {
		sb.shutDown()
	}

	if vtrace.GlobalTracer != nil {
		vtrace.GlobalTracer.Close()
	}

}

func (sb *ServiceBase) initRegisterEngines() error {

	serverConfig := new(vconfig.ServerConfig)
	err := sb.center.UnmarshalWithNameSpace(vconfig.Server, parse.PropertiesTagName, serverConfig)
	if err != nil {
		return err
	}

	regEngine, err := register.GetRegisterEngine(&common.RegisterConfig{
		RegistrationType: common.RegistrationType(serverConfig.RegisterType),
		Cluster:          serverConfig.RegisterCluster,
	})
	if err != nil {
		return err
	}
	sb.register = regEngine

	if serverConfig.NeedMetric {
		metricEngine, err := consul.NewRegistry(&consul.RegisterConfig{
			Host:  serverConfig.ConsulHost,
			Port:  serverConfig.ConsulPort,
			Token: serverConfig.ConsulToken,
		})
		if err == nil {
			sb.metricRegister = metricEngine
		}
	}

	logLevel, err := zapcore.ParseLevel(serverConfig.LogLevel)
	if err != nil {
		logLevel = zapcore.WarnLevel
	}
	vlog.InitLogger(sb.buildLogPath(), vlog.LogFile, logLevel, vlog.JsonFormatType)

	return nil
}

func (sb *ServiceBase) initCenter(args *servArgs) error {

	cfg, err := sb.parseConfigEnv(args)
	if err != nil {
		return err
	}

	center, err := vconfig.NewCenter(cfg)
	if err != nil {
		return err
	}

	sb.center = center
	return nil
}

func (sb *ServiceBase) parseConfigEnv(args *servArgs) (*vconfig.CenterConfig, error) {
	centerConfig, err := vconfig.ParseConfigEnv()
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(centerConfig.Port, 10, 64)
	if err != nil {
		return nil, err
	}

	return &vconfig.CenterConfig{
		AppID:   args.serviceGroup + common.Slash + args.serviceName,
		Cluster: centerConfig.Cluster,
		// todo db、mq、redis config
		Namespace:        []string{vconfig.Application, vconfig.Server, vconfig.ServerDB},
		IP:               centerConfig.IP,
		Port:             int(port),
		IsBackupConfig:   centerConfig.IsBackupConfig,
		BackupConfigPath: sb.buildCenterBackupPath(),
		MustStart:        centerConfig.MustStart,
	}, nil
}

func (sb *ServiceBase) buildCenterBackupPath() string {
	return fmt.Sprintf("%s/%s/%s/%d/%s", common.TmpPath, sb.group, sb.name+sb.lane, sb.version, vconfig.BackupPath)
}

func (sb *ServiceBase) buildLogPath() string {
	return fmt.Sprintf("%s/%s/%s/%d/%s", common.TmpPath, sb.group, sb.name+sb.lane, sb.version, vlog.LogPath)
}

func (sb *ServiceBase) initMetric(ctx context.Context) error {
	if sb.metricRegister == nil {
		return nil
	}

	metric.InitBaseMetric(ctx, sb.group, sb.name, sb.ID)
	serv, err := sb.powerServices(ctx, map[common.ServiceType]common.Processor{
		common.Metric: &vcollector.MetricProcessor{},
	})
	if err != nil {
		return err
	}
	serviceInfo, ok := serv.ServList[common.Metric]
	if !ok {
		return nil
	}

	_, err = sb.metricRegister.Register(ctx, sb.FullServiceRegisterPath(), serviceInfo.Addr, common.DefaultTTl)
	if err != nil {
		return err
	}

	return nil
}

func (sb *ServiceBase) powerServices(ctx context.Context, props map[common.ServiceType]common.Processor) (*common.RegisterServiceInfo, error) {
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

	return &common.RegisterServiceInfo{
		Lane:     sb.lane,
		ServList: serv,
	}, nil
}
