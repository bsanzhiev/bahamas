package customer_kafka

import (
	"context"
	"customers_service/internal/infrastructure/config"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

// NewPublisher creates a new Kafka publisher
// Return directly the Publisher struct
// NewWriiter are deprecated
func NewPublisher(cfg config.KafkaConfig) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(cfg.Brokers...),
			Topic: cfg.Topic,
		},
	}
}

func (p *Publisher) Publish(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Message published successfully: key: %s, value: %s", key, string(value))
	return nil
}

func (p *Publisher) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka writer: %w", err)
	}
	return nil
}
