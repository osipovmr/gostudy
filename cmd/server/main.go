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
	logger.InitLogger()

	cfg := config.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application := app.NewApp(cfg)

	if err := application.Run(ctx); err != nil {
		slog.Error("app stopped with error", "error", err)
	}

	slog.Info("app exited")
}
