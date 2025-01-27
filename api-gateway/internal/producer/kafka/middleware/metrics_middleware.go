package middleware

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
	consumer "todo/analytics-service/domain"
	"todo/pkg/http/handlers"
	"todo/pkg/infra/broker/producer"
	pkglog "todo/pkg/log"
)

func WithMetricsProducer(ctx context.Context, logger *slog.Logger, producer producer.Producer) handlers.RouterOption {
	return func(r chi.Router) {
		r.Use(NewMetricsProducerMiddleware(ctx, logger, producer))
	}
}

func NewMetricsProducerMiddleware(ctx context.Context, log *slog.Logger, producer producer.Producer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/producer"),
		)

		log.Info("Producer middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				duration := time.Since(t1)

				metrics := consumer.ResponseMetrics{
					Method:     r.Method,
					Path:       r.URL.Path,
					StatusCode: ww.Status(),
					Time:       duration,
					Size:       ww.BytesWritten(),
				}

				go func() {
					if err := producer.Produce(ctx, metrics); err != nil {
						log.Error("failed to produce metrics", pkglog.Err(err))
					}
				}()
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
