package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log/slog"
	_ "todo/api-gateway/docs"
	"todo/api-gateway/internal/api/http"
	"todo/api-gateway/internal/clients/todo/grpc"
	"todo/api-gateway/internal/config"
	pkgconfig "todo/pkg/config"
	"todo/pkg/http/handlers"
	"todo/pkg/http/server"
	pkglog "todo/pkg/log"
	"todo/pkg/shutdown"
)

//	@title			Task Manager API Gateway
//	@version		1.0
//	@description	API Gateway for GRPC service of task management
//	@termsOfService	http://swagger.io/terms/

//	@host		localhost:8080
//	@BasePath	/api/v1/

const ConfigEnvVar = "API_GATEWAY_CONFIG"
const APIPath = "/api/v1/"

func main() {
	cfg := pkgconfig.MustLoad[config.Config](ConfigEnvVar)

	log, file := pkglog.NewLogger(cfg.Logger)
	defer func() { _ = file.Close() }()
	slog.SetDefault(log)
	log.Info("Starting API-Gateway", slog.Any("config", cfg))

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return shutdown.ListenSignal(ctx, log)
	})

	grpcClient, err := grpc.New(
		ctx,
		log,
		cfg.GRPCClient.Address,
		cfg.GRPCClient.Timeout,
		cfg.GRPCClient.Retries,
	)
	if err != nil {
		pkglog.Fatal(log, "error while setting new grpc client: ", err)
	}

	taskHandler := http.NewTaskHandler(log, grpcClient)

	publicHandler := handlers.NewHandler(APIPath,
		handlers.WithLogging(log),
		handlers.WithRecover(),
		handlers.WithProfilerHandlers(),
		handlers.WithSwagger(),
		handlers.WithErrHandlers(),
		handlers.WithHealthHandler(),
		taskHandler.WithTaskHandlers(),
	)

	g.Go(func() error {
		return server.RunServer(ctx, cfg.HTTPServer.Address,
			publicHandler,
			cfg.HTTPServer.WriteTimeout,
			cfg.HTTPServer.ReadTimeout,
			cfg.HTTPServer.IdleTimeout)
	})

	err = g.Wait()
	if err != nil && !errors.Is(err, errors.New("operating system signal")) {
		log.Error("Exit reason", slog.String("error", err.Error()))
	}
}
