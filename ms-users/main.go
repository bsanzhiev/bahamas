package main

import (
	"context"
	"encoding/json"
	"fmt"

	"log"
	"os"

	gatewayTypes "github.com/bsanzhiev/bahamas/ms-gateway/types"
	"github.com/bsanzhiev/bahamas/ms-users/migrations"
	"github.com/gofiber/fiber/v2"

	"github.com/IBM/sarama"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// TODO: 1. Remove HTTP Routes +
// TODO: 2. Set Up Kafka Consumer
// TODO: 3. Message Processing Logic
// TODO: 4.Response Handling

func main() {
	ctx := context.Background()

	// Connect to Kafka brokers
	brokers := []string{"localhost:9092"}
	//consumerGroup := "users_consumer_group"
	//topics := []string{"users_requests"}
	topics := "users_requests"

	// Initialize Kafka consumer
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create a new consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v: ", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topics, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to get partition consumer: %v", err)
	}

	// Create a new synchronous producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Failed to close partition consumer: %v", err)
		}
		if err := consumer.Close(); err != nil {
			log.Printf("Failed to close consumer: %v", err)
		}
	}()

	go func() {
		for {
			select {
			case err := <-partitionConsumer.Errors():
				log.Printf("Error: %v", err)
			case msg := <-partitionConsumer.Messages():
				// Process incoming messages
				var requestData = gatewayTypes.RequestData{}
				err := json.Unmarshal(msg.Value, &requestData)
				if err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}

				// Extract action and request data
				action := requestData.Action
				data := requestData.Data

				// Perform corresponding operations based on action
				switch action {
				case "user_list":
					// handle get all users
					fmt.Printf("Data: %v", data)
				case "user_by_id":
					// handle get user by id
					fmt.Printf("User ID: %v", data)
				default:
					log.Printf("Unknown action: %v", action)
				}

				// Generate response
				var responseData = gatewayTypes.ResponseData{}
				responseData.Status = 200
				responseData.Message = "Success"
				responseData.Data = "Response Data"

				// Send response to Kafka topic ('users_responses')
				responseTopic := "users_responses"
				responseJSON, err := json.Marshal(responseData)
				if err != nil {
					log.Printf("Failed to marshal response data: %v", err)
					continue
				}
				producerMsg := sarama.ProducerMessage{
					Topic: responseTopic,
					Value: sarama.ByteEncoder(responseJSON),
				}

				if _, _, err := producer.SendMessage(&producerMsg); err != nil {
					log.Printf("Failed to send response message: %v", err)
					continue
				}
				//consumer.MarkMessage(msg, "")
			}
		}
	}()

	// Main Users app =================================
	// Получаем строку подключения
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatalf("Error loading .env file: %v", errEnv)
	}
	urlDB := os.Getenv("CONNECTION_STRING")
	// Создаем пул подключений к базе данных
	dbPool, errPool := pgxpool.New(ctx, urlDB)
	if errPool != nil {
		log.Fatalf("Failed to create pool: %v", errPool)
	}
	defer dbPool.Close()

	fmt.Println("Successfully connected to database!")

	// Делаем миграцию
	migrator, err := migrations.NewMigrator(ctx, dbPool)
	if err != nil {
		panic(err)
	}

	// get the current migration status
	now, exp, info, err := migrator.Info()
	if err != nil {
		panic(err)
	}
	if now < exp {
		// migration is required, dump out the current state
		// and perform the migration
		println("migration needed, current state:")
		println(info)

		err = migrator.Migrate(ctx)
		if err != nil {
			panic(err)
		}
		println("migration successful!")
	} else {
		println("no database migration needed")
	}

	app := fiber.New(
		fiber.Config{
			AppName: "Bahamas Users Service",
		},
	)

	app.Get("/alive", Alive)

	// Routing using REST API
	// Provide connection pool and context
	//userController := &controllers.UserController{
	//	DBPool: dbPool,
	//	Ctx:    ctx,
	//}
	// Grouping routes
	//api := app.Group("/api/v1")
	//api.Get("/users", func(c *fiber.Ctx) error {
	//	return c.SendString("Return all accounts v1")
	//})
	//controllers.RegisterUserRoutes(api, userController)

	if err := app.Listen(":7002"); err != nil {
		fmt.Printf("Error starting User server: %s\n", err)
	}
	// Main app =================================================
}

// Alive Readiness Check
func Alive(c *fiber.Ctx) error {
	defer func() {
		err := c.JSON(fiber.Map{"alive": true, "ready": true, "service": "users"})
		if err != nil {
			return
		}
	}()
	return nil
}
