package message_consumer

import (
	"context"
	kafka_handlers "customers_service/internal/interfaces/message_consumer/handlers"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaCosumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer create new instance of Kafka consumer
func NewKafkaConsumer(brockers []string, topic string) *KafkaCosumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brockers,
		Topic:   topic,
		GroupID: "customer-service-group",
	})

	return &KafkaCosumer{
		reader: reader,
	}
}

func (c *KafkaCosumer) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer...")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			log.Printf("Recieved message: key: %s, value: %s", string(msg.Key), string(msg.Value))
			kafka_handlers.HandleLoanApproved(ctx, msg.Value)
		}
	}
}

func (c *KafkaCosumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka reader: %w", err)
	}
	return nil
}
