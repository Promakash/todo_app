package service

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"todo/metrics-service/domain"
	"todo/metrics-service/internal/usecases"
	"todo/pkg/infra/broker/consumer"
	pkglog "todo/pkg/log"
)

type MetricsService struct {
	log             *slog.Logger
	consumer        consumer.Consumer
	requestsTotal   *prometheus.CounterVec
	responseSize    *prometheus.HistogramVec
	requestDuration *prometheus.HistogramVec
}

func NewMetricsService(log *slog.Logger, consumer consumer.Consumer) usecases.MetricsScrapper {
	reqTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	respSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"method", "path"},
	)

	reqDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	prometheus.MustRegister(reqTotal, respSize, reqDuration)

	return &MetricsService{
		log:             log,
		consumer:        consumer,
		requestsTotal:   reqTotal,
		responseSize:    respSize,
		requestDuration: reqDuration,
	}
}

func (s *MetricsService) ScrapeMetrics() {
	const op = "MetricsService.ScrapeMetrics"
	log := slog.With(slog.String("op", op))
	for {
		var metric domain.ResponseMetrics
		err := s.consumer.Consume(&metric)
		if err != nil {
			log.Error("Error while consuming from kafka", pkglog.Err(err))
			return
		}
		s.requestsTotal.WithLabelValues(metric.Method, metric.Path, fmt.Sprintf("%d", metric.StatusCode)).Inc()
		s.responseSize.WithLabelValues(metric.Method, metric.Path).Observe(float64(metric.Size))
		s.requestDuration.WithLabelValues(metric.Method, metric.Path).Observe(float64(metric.Time.Milliseconds()))
	}
}
