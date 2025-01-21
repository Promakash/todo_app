package shutdown

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func ListenSignal(ctx context.Context, logger *slog.Logger) {
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ctx.Done():
		return
	case sig := <-sigquit:
		logger.Info("Captured signal: ", sig)
		return
	}
}
