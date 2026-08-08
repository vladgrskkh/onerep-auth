package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladgrskkh/onerep-auth/internal/application"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app, err := application.New(application.WithLogger(logger))
	if err != nil {
		logger.Error("failed to create app", "error", err)
		return err
	}

	go func() {
		if err := app.Run(ctx); err != nil {
			logger.Error("failed to run", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), app.ShutdownTimeout())
	defer shutdownCancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		logger.Error("force shutdown", "error", err)
		return err
	}

	logger.Info("stopped")
	return nil
}
