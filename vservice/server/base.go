package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Snowlights/tool/cache/vredis"
	"github.com/Snowlights/tool/parse"
	"github.com/Snowlights/tool/vconfig"
	"github.com/Snowlights/tool/vlog"
	"github.com/Snowlights/tool/vmongo"
	"github.com/Snowlights/tool/vmq"
	"github.com/Snowlights/tool/vprometheus/metric"
	"github.com/Snowlights/tool/vprometheus/vcollector"
	"github.com/Snowlights/tool/vservice/common"
	"github.com/Snowlights/tool/vservice/server/engine"
	"github.com/Snowlights/tool/vservice/server/register"
	"github.com/Snowlights/tool/vservice/server/register/consul"
	"github.com/Snowlights/tool/vsql"
	"github.com/Snowlights/tool/vtrace"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

type ServiceBase struct {
	center vconfig.Center

	// todo: redis can refactor like mongo and vsql
	redisClient  *vredis.RedisClient
	mongoManager *vmongo.Manager
	mqManager    *vmq.Manager
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

	err = servBase.initRedis(ctx)
	if err != nil {
		vlog.ErrorF(ctx, "init redis error: %v", err)
		err = nil
	}

	err = servBase.initMQ(ctx)
	if err != nil {
		vlog.ErrorF(ctx, "init mq error: %v", err)
		err = nil
	}

	err = servBase.initMongo(ctx)
	if err != nil {
		vlog.ErrorF(ctx, "init mongo error: %v", err)
		err = nil
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

func (sb *ServiceBase) GetRedisClient(ctx context.Context) *vredis.RedisClient {
	return sb.redisClient
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

	if sb.mongoManager != nil {
		sb.mongoManager.Close()
	}

	if sb.mqManager != nil {
		sb.mqManager.Close()
	}

	if vtrace.GlobalTracer != nil {
		vtrace.GlobalTracer.Close()
	}

}

func (sb *ServiceBase) initRedis(ctx context.Context) error {
	redisConfig := new(vredis.RedisConfig)
	err := sb.center.UnmarshalWithNameSpace(vconfig.Redis, parse.PropertiesTagName, redisConfig)
	if err != nil {
		return err
	}

	redisClient, err := vredis.NewRedisClient(ctx, redisConfig)
	if err != nil {
		return err
	}
	sb.redisClient = redisClient
	return nil
}

func (sb *ServiceBase) initMongo(ctx context.Context) error {

	mongoManager, err := vmongo.NewManager(sb.center)
	if err != nil {
		return err
	}
	sb.mongoManager = mongoManager
	return nil
}

func (sb *ServiceBase) initMQ(ctx context.Context) error {

	mqManager, err := vmq.NewManager(sb.center)
	if err != nil {
		return err
	}
	sb.mqManager = mqManager
	return nil
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
