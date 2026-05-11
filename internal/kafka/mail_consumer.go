package kafka

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string
	GroupID string
}

type MailConsumer struct {
	cfg   Config
	topic string
}

func NewMailConsumer(cfg Config, topic string) *MailConsumer {
	return &MailConsumer{
		cfg:   cfg,
		topic: topic,
	}
}

func (c *MailConsumer) Run(ctx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.cfg.Brokers,
		GroupID: c.cfg.GroupID,
		Topic:   c.topic,
	})
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			slog.Error("failed to close connection", "error", err)
		}
	}(reader)

	slog.Info("mail consumer started", "topic", c.topic)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		slog.Info("consumer",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"value", string(msg.Value),
		)
	}
}
