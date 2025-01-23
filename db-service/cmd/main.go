package main

import (
	"context"
	"log/slog"
	"todo/db-service/internal/app/grpc"
	"todo/db-service/internal/config"
	"todo/db-service/internal/repository/postgres"
	"todo/db-service/internal/service"
	pkgconfig "todo/pkg/config"
	"todo/pkg/infra"
	"todo/pkg/infra/cache/redis"
	pkglog "todo/pkg/log"
	"todo/pkg/shutdown"
)

const ConfigEnvVar = "DB_SERVICE_CONFIG"

func main() {
	cfg := pkgconfig.MustLoad[config.Config](ConfigEnvVar)

	log := pkglog.NewLogger("debug", "json")
	slog.SetDefault(log)
	log.Info("Starting dbService", slog.Any("config", cfg))

	dbPool, err := infra.NewPostgresPool(cfg.PG)
	if err != nil {
		pkglog.Fatal(log, "error while setting new postgres connection: ", err)
	}
	defer dbPool.Close()
	taskRepo := postgres.NewTaskRepository(dbPool)

	redisClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		pkglog.Fatal(log, "error while setting new redis connection: ", err)
	}
	defer redis.ShutdownClient(redisClient)
	cacheService := redis.NewRedisService(redisClient)

	taskService := service.NewTaskService(log, taskRepo, cacheService, cfg.Redis.TTL)

	application := grpc.NewApp(log, cfg.GRPC.Port, taskService)

	application.MustRun()

	shutdown.ListenSignal(context.Background(), log)

	application.Stop()

	log.Info("application stopped")
}
