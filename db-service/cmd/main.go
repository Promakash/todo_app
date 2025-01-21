package main

import (
	"context"
	"log/slog"
	"os"
	"todo/db-service/cmd/grpc"
	sConfig "todo/db-service/internal/config"
	"todo/db-service/internal/repository/postgres"
	"todo/db-service/internal/service"
	"todo/pkg/config"
	"todo/pkg/infra"
	pkglog "todo/pkg/log"
	"todo/pkg/shutdown"
)

const ConfigEnvVar = "DB_SERVICE_CONFIG"

func main() {
	cfg := config.MustLoad[sConfig.Config](ConfigEnvVar)

	log := pkglog.NewLogger("debug", "json")
	slog.SetDefault(log)
	log.Info("Starting dbService", slog.Any("config", cfg))

	dbPool, err := infra.NewPostgresPool(cfg.PG)
	if err != nil {
		log.Error("error while setting new postgres connection: ", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	taskRepo := postgres.NewTaskRepository(dbPool)
	taskService := service.NewTaskService(taskRepo, log)

	application := grpc.NewApp(log, cfg.GRPC.Port, taskService)

	application.MustRun()

	shutdown.ListenSignal(context.Background(), log)

	application.Stop()

	log.Info("application stopped")
}
