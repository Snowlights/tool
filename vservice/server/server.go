package server

import (
	"context"
	"flag"
	"github.com/Snowlights/tool/vlog"
	"github.com/Snowlights/tool/vservice/common"
	"os"
	"os/signal"
	"syscall"
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
	serviceLane  string

	version int64
}

func GetServBase() common.ServerBase {
	return server.serviceBase
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

	err = s.serviceBase.Register(ctx, props)
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
				return
			}
		}
	}
}

func (s *Server) parseServiceInfo() (*servArgs, error) {
	var serv, group, lane string
	var version int64

	flag.StringVar(&serv, "serv", "censor", "service name")
	flag.StringVar(&group, "group", "base/talent", "service group")
	flag.StringVar(&lane, "lane", "", "service lane")
	flag.Int64Var(&version, "version", 1, "service build version")

	flag.Parse()

	if len(serv) == 0 {
		return nil, common.ServiceNameIsNil
	}

	if len(group) == 0 {
		return nil, common.ServiceGroupIsNil
	}

	return &servArgs{
		serviceName:  serv,
		serviceGroup: group,
		version:      version,
		serviceLane:  lane,
	}, nil
}
