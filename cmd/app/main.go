package main

import (
	"context"
	"log/slog"
	"time"

	"os/signal"
	"syscall"

	"gostudy/internal/app"
	"gostudy/internal/config"

	"gostudy/internal/logger"

	"github.com/segmentio/kafka-go"
)

func main() {
	logger.Init()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(cfg)
	if err != nil {
		slog.Error("failed to init app", "error", err)
		return
	}

	dialer := &kafka.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", cfg.KAFKAAddr)
	if err != nil {
		slog.Error(err.Error())
	}
	defer conn.Close()

	err = conn.CreateTopics(
		kafka.TopicConfig{
			Topic:             "__consumer_offsets",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
	if err != nil {
		slog.Error(err.Error())
	}
	err = conn.CreateTopics(
		kafka.TopicConfig{
			Topic:             "registration_topic",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info("topic created")

	if err := application.Run(ctx); err != nil {
		slog.Error("app stopped with error", "error", err)
	}

	slog.Info("app exited")
}
