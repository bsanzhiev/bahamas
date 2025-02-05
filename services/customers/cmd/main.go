package main

import (
	"log"

	"github.com/bsanzhiev/bahamas/services/customers/internal/infrastructure/config"
	customer_kafka "github.com/bsanzhiev/bahamas/services/customers/internal/infrastructure/kafka"
	"github.com/bsanzhiev/bahamas/services/customers/internal/infrastructure/postgres"
	"github.com/bsanzhiev/bahamas/services/customers/internal/interfaces/message_consumer"

	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Init repositories
	pgRepo := postgres.NewRepository(cfg.Postgres)
	defer pgRepo.Close()

	// Kafka producer
	kafkaProducer := customer_kafka.NewPublisher(cfg.Kafka)

	// Use cases
	createUC := usecases.NewCreateCustomerUseCase(pgRepo, kafkaProducer)

	// gRPC server
	grpcHandler := grpc_handler.NewCustomerGRPCHandler(createUC)
	grpcServer := grpc.NewServer(grpcHandler)

	// Start Kafka event consumer
	kafkaConsumer := message_consumer.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	kafkaConsumer.Subscribe("loan_approved", loadApprovedHandler)

	// Start servers
	go grpcServer.Start(cfg.GRPC.Port)
	kafkaConsumer.Run()

}
