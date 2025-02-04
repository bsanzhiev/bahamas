package main

import (
	"customers_service/internal/infrastructure/config"
	customer_kafka "customers_service/internal/infrastructure/kafka"
	"customers_service/internal/infrastructure/postgres"
	"customers_service/internal/interfaces/message_consumer"
	"log"

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
