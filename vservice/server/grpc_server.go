package server

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"vtool/vprometheus/metric"
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
		monitorServerInterceptor(),
	)

	streamInterceptors = append(streamInterceptors,
		grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		monitorStreamServerInterceptor(),
	)

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)))
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)))

	serv := grpc.NewServer(opts...)
	return serv
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
