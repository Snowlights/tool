package server

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"vtool/vlog"
	"vtool/vservice/common"
)

var server *Server

func init() {
	server = &Server{}
}

type Server struct {
	serviceBase common.ServerBase
}

type servArgs struct {
	serviceName  string
	serviceGroup string

	logDir       string
	registerType common.RegistrationType
}

func ServService(props map[common.ServiceType]common.Processor) error {
	return server.serv(props)
}

func (s *Server) serv(props map[common.ServiceType]common.Processor) error {
	ctx := context.Background()
	args, err := s.parseServiceInfo()
	if err != nil {
		return err
	}

	servBase, err := NewServiceBase(ctx, args)
	if err != nil {
		return err
	}

	s.serviceBase = servBase

	err = s.serverIns(ctx, props)
	if err != nil {
		return err
	}

	s.awaitSignal()
	return nil
}

func (s *Server) awaitSignal() {
	c := make(chan os.Signal, 1)
	ctx := context.Background()
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGPIPE}
	signal.Reset(signals...)
	signal.Notify(c, signals...)

	for {
		select {
		case signalInfo := <-c:
			vlog.InfoF(ctx, "receive a signal:%s", signalInfo.String())

			if signalInfo.String() == syscall.SIGTERM.String() {
				vlog.InfoF(ctx, "receive a signal: %s, stop service", signalInfo.String())
				s.serviceBase.Stop()
				<-(chan int)(nil)
			}
		}
	}
}

func (s *Server) parseServiceInfo() (*servArgs, error) {
	var serv, logDir, group string
	var registerType int64
	flag.StringVar(&serv, "serv", "", "service name")
	flag.StringVar(&logDir, "logDir", "", "service log dir")
	flag.StringVar(&group, "group", "", "service group")
	flag.Int64Var(&registerType, "regType", 0, "service register type")
	flag.Parse()

	if len(serv) == 0 {
		return nil, common.ServiceNameIsNil
	}

	if len(group) == 0 {
		return nil, common.ServiceGroupIsNil
	}

	if registerType == 0 {
		return nil, common.RegisterTypeIsNil
	}

	if len(logDir) == 0 {
		return nil, common.LogDirIsNil
	}

	return &servArgs{
		serviceName:  serv,
		serviceGroup: group,
		logDir:       logDir,
		registerType: common.RegistrationType(registerType),
	}, nil
}

func (s *Server) serverIns(ctx context.Context, props map[common.ServiceType]common.Processor) error {
	return s.serviceBase.Register(ctx, props)
}
