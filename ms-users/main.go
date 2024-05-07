package main

import (
	"context"
	"encoding/json"
	"fmt"

	"log"
	"os"

	// gatewayTypes "github.com/bsanzhiev/bahamas/ms-gateway/types"
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
	consumerGroup := "users_consumer_group"

	// Init Kafka consumer
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, consumerGroup, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v: ", err)
	}
	defer func(consumer sarama.ConsumerGroup) {
		err := consumer.Close()
		if err != nil {

		}
	}(consumer)

	topics := []string{"users_requests"}
	err = consumer.Consume(ctx, topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %v", err)
	}

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

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Printf("Error: %v", err)
			case msg := <-consumer.Messages():
				// Process incoming messages
				var requestData = gateway.RequestData{}
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
				default:
					log.Printf("Unknown action: %v", action)
				}

				// Generate response
				var responseData = ResponseData
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
			}
		}
	}()
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
