package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"regexp"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/SangBejoo/Template/init/config"
	"github.com/SangBejoo/Template/init/infra"
	"github.com/SangBejoo/Template/internal/dependency"
)

func RunGRPCServer(ctx context.Context, cfg config.Main, repo infra.Repository) (*grpc.Server, error) {
	grpcPort := fmt.Sprintf(":%d", cfg.GrpcServer.Port)
	grpcConn, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen port: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metadataInterceptor),
		grpc.ChainUnaryInterceptor(
			apmgrpc.NewUnaryServerInterceptor(),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(grpcRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			apmgrpc.NewStreamServerInterceptor(),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(grpcRecoveryHandler)),
		),
	)

	dependency.InitGrpcDependency(grpcServer, repo)
	reflection.Register(grpcServer)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	for name := range grpcServer.GetServiceInfo() {
		healthServer.SetServingStatus(name, healthpb.HealthCheckResponse_SERVING)
	}
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go grpcServer.Serve(grpcConn)
	slog.Info(fmt.Sprintf("server grpc listening at %v", grpcConn.Addr()))
	return grpcServer, nil
}

func grpcRecoveryHandler(ctx context.Context, panic interface{}) error {
	newLineRegex := regexp.MustCompile(`\r?\n`)
	stackTrace := newLineRegex.ReplaceAllString(string(debug.Stack()), " ")
	slog.Error("panic happened",
		"panic_message", panic,
		"panic_stacktrace", stackTrace)
	return status.Errorf(codes.Internal, "server error happened")
}

func metadataInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	newCtx := metadata.NewIncomingContext(ctx, md)
	return handler(newCtx, req)
}
