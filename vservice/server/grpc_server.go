package server

import (
	"context"
	opentrace_go_grpc "github.com/Snowlights/gogrpc"
	"github.com/Snowlights/tool/vprometheus/metric"
	"github.com/Snowlights/tool/vservice/common"
	"github.com/Snowlights/tool/vtrace"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func NewGrpcServerWithInterceptor() *grpc.Server {
	return buildServer()
}

func buildServer() *grpc.Server {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	// todo : add trace, rate limit interceptor
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}
	unaryInterceptors = append(unaryInterceptors,
		grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		opentrace_go_grpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), opentrace_go_grpc.SpanDecorator(SpanDecorator)),
		monitorServerInterceptor(),
	)

	streamInterceptors = append(streamInterceptors,
		grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		opentrace_go_grpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer(), opentrace_go_grpc.SpanDecorator(SpanDecorator)),
		monitorStreamServerInterceptor(),
	)

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)))
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)))

	serv := grpc.NewServer(opts...)
	return serv
}

func SpanDecorator(ctx context.Context,
	span opentracing.Span,
	method string,
	req, resp interface{},
	grpcError error) {
	// todo more tag info added
	span.SetTag("grpc.method", method)
	servBase := GetServBase()
	if servBase != nil {
		servInfo := servBase.ServInfo()
		if servInfo != nil {
			span.SetTag(vtrace.Lane, servInfo.Lane)
			span.SetTag(vtrace.ServType, common.Grpc)
			serv, ok := servInfo.ServList[common.Grpc]
			if ok {
				span.SetTag(vtrace.ServIP, serv.Addr)
				span.SetTag(vtrace.EngineType, serv.Type)
			}
		}
	}

}

func monitorServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		st := time.Now()
		resp, err = handler(ctx, req)
		metric.StatApi(info.FullMethod, time.Since(st))
		return
	}
}

func monitorStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		st := time.Now()
		err = handler(srv, ss)
		metric.StatApi(info.FullMethod, time.Since(st))
		return
	}
}

func recoveryFunc(p interface{}) (err error) {
	return status.Errorf(codes.Unknown, "grpc server recoveryFunc %v", p)
}
