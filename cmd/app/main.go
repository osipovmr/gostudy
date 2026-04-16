package main

import (
	"context"
	"log/slog"

	"os/signal"
	"syscall"

	"gostudy/internal/app"
	"gostudy/internal/config"

	"gostudy/internal/logger"
)

func main() {
	logger.Init()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application := app.New(cfg)

	if err := application.Run(ctx); err != nil {
		slog.Error("app stopped with error", "error", err)
	}

	slog.Info("app exited")
}
