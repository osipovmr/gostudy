package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer                *kafka.Writer
	mailRegistrationTopic string
}

var (
	ErrKafkaProducer = errors.New("fail to send message")
)

type Producer interface {
	SendRegistrationMessage(ctx context.Context, userEmail string) error
}

func NewProducer(brokers []string,
	mailRegistrationTopic string,
) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
		mailRegistrationTopic: mailRegistrationTopic,
	}

}

func (p *producer) SendRegistrationMessage(ctx context.Context, userEmail string) error {
	err := p.send(ctx, p.mailRegistrationTopic, nil, []byte(userEmail))
	if err != nil {
		return ErrKafkaProducer
	}
	slog.Info("producer: topic %s, message: %s", p.mailRegistrationTopic, userEmail)
	return nil
}

func (p *producer) send(ctx context.Context, topic string, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}
