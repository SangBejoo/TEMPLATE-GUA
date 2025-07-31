package main

import (
	"context"
	"fmt"
	"log/slog"
	_ "net/http/pprof"
	"time"

	"github.com/SangBejoo/Template/init/config"
	"github.com/SangBejoo/Template/init/infra"
	"github.com/SangBejoo/Template/init/logger"
	"github.com/SangBejoo/Template/init/server"
	"github.com/SangBejoo/Template/util"
)

var cfg *config.Main

func init() {
	cfg = config.Load()
	logger.Load(*cfg)
}

func main() {
	// TODO: complete usecase implementation in usecase folder
	repo := infra.LoadRepository(*cfg)
	defer func() {
		if errClose := repo.Close(); errClose != nil {
			slog.Error(fmt.Sprintf("failed to close repositories: %v", errClose))
		}
	}()

	ctx := context.Background()
	grpcServer, err := server.RunGRPCServer(ctx, *cfg, *repo)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to run grpc server: %v", err))
	}

	err = server.RunGatewayRestServer(ctx, *cfg, *repo)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to run gateway rest server: %v", err))
	}

	wait := util.GracefulShutdown(ctx, 5*time.Second, map[string]util.Operation{
		"grpc": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
	<-wait
}
