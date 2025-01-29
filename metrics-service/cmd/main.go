package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
	"todo/metrics-service/internal/config"
	"todo/metrics-service/internal/service"
	pkgconfig "todo/pkg/config"
	"todo/pkg/http/handlers"
	"todo/pkg/http/server"
	kafkaconsumer "todo/pkg/infra/broker/consumer/kafka"
	pkglog "todo/pkg/log"
	"todo/pkg/shutdown"
)

const ConfigEnvVar = "METRICS_SERVICE_CONFIG"

func main() {
	cfg := pkgconfig.MustLoad[config.Config](ConfigEnvVar)

	log, _ := pkglog.NewLogger(cfg.Logger)
	slog.SetDefault(log)
	log.Info("Starting Metrics service", slog.Any("config", cfg))

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return shutdown.ListenSignal(ctx, log)
	})

	kafkaCons, err := kafkaconsumer.NewConsumer(cfg.KafkaConsumer)
	if err != nil {
		pkglog.Fatal(log, "error while setting new kafka consumer: ", err)
	}
	defer kafkaCons.Close()

	metricsService := service.NewMetricsService(log, kafkaCons)
	go func() {
		metricsService.ScrapeMetrics()
	}()

	publicHandler := handlers.NewHandler("/",
		handlers.WithRecover(),
		handlers.WithMetricsHandler(),
		handlers.WithErrHandlers(),
	)

	g.Go(func() error {
		return server.RunServer(ctx, cfg.HTTPServer.Address,
			publicHandler,
			time.Minute,
			time.Minute,
			time.Minute)
	})

	err = g.Wait()
	if err != nil && !errors.Is(err, errors.New("operating system signal")) {
		log.Error("Exit reason", slog.String("error", err.Error()))
	}
}
