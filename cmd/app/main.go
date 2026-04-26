package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
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

	if err := createTopicsWithRetry(ctx, cfg.KAFKAAddr, []kafka.TopicConfig{
		{
			Topic:             "__consumer_offsets",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		{
			Topic:             "registration_topic",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}, 2, 10*time.Second); err != nil {
		slog.Error("failed to create topics", "error", err)
		return
	}

	slog.Info("topics created")

	if err := application.Run(ctx); err != nil {
		slog.Error("app stopped with error", "error", err)
	}

	slog.Info("app exited")
}

func createTopicsWithRetry(
	ctx context.Context,
	brokerAddr string,
	topics []kafka.TopicConfig,
	attempts int,
	delay time.Duration,
) error {
	var lastErr error

	for i := 1; i <= attempts; i++ {
		if err := createTopics(ctx, brokerAddr, topics); err != nil {
			lastErr = err
			slog.Error("topic creation failed", "attempt", i, "error", err)

			if i < attempts {
				timer := time.NewTimer(delay)
				select {
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()
				case <-timer.C:
				}
			}
			continue
		}

		return nil
	}

	return fmt.Errorf("create topics failed after %d attempts: %w", attempts, lastErr)
}

func createTopics(ctx context.Context, brokerAddr string, topics []kafka.TopicConfig) error {
	dialer := &kafka.Dialer{Timeout: 10 * time.Second}

	conn, err := dialer.DialContext(ctx, "tcp", brokerAddr)
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))

	controllerConn, err := dialer.DialContext(ctx, "tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()

	if err := controllerConn.CreateTopics(topics...); err != nil {
		return fmt.Errorf("create topics: %w", err)
	}

	return nil
}
