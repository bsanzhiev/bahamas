package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Создаем новый экземпляр Fiber
	app := fiber.New(
		fiber.Config{
			AppName: "Bahama API Gateway",
		},
	)

	// Проверка состояния работы сервиса
	app.Get("/alive", Alive)

	// Проверка аутентификации
	// Middleware for check auth token
	// app.Use(func(c *fiber.Ctx) error {
	// 	token := c.Get("Authorization")

	// 	if token == "" {
	// 		c.Status(401).JSON(fiber.Map{"error": "No valid token"})
	// 		return c.Redirect("/login")
	// 	}
	// 	return c.Next()
	// })

	// Разбираем полученный запрос;
	// Узнаем имя сервиса - сопоставляем с адресами конкретного сервиса.

	app.Post("/", HandleRequest)

	// Запуск шлюза
	go func() {
		if err := app.Listen(":7001"); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()
	select {}
}

// Alive Проверка работы
func Alive(c *fiber.Ctx) error {
	defer func() {
		err := c.JSON(fiber.Map{"alive": true, "ready": true, "service": "gateway"})
		if err != nil {
			return
		}
	}()
	return nil
}

// RequestData Тип для запросов
type RequestData struct {
	Service string                 `json:"service"` // Имя сервиса
	Action  string                 `json:"action"`  // Имя операции
	Data    map[string]interface{} `json:"data"`    // Объект запроса
}

// ResponseData Тип для ответов
type ResponseData struct {
	Status  int         `json:"status"`  // Код ответа
	Message string      `json:"message"` // Сообщение об ошибке или результате
	Data    interface{} `json:"data"`    // Объект ответа
}

// HandleRequest Обработка входящих запросов
func HandleRequest(c *fiber.Ctx) error {

	rawBody := c.Body()
	log.Printf("Raw request body: %s\n", rawBody)

	requestData := RequestData{}
	if err := c.BodyParser(&requestData); err != nil {
		return err
	}
	log.Printf("Request data from body: %v\n", requestData)

	// Send data to users_requests topic
	errSend := SendToKafka(&requestData)
	if errSend != nil {
		return c.Status(500).SendString(errSend.Error())
	}

	// Get response data from topic
	responseData := ResponseData{}
	errGet := GetFromKafka(&requestData, &responseData)
	if errGet != nil {
		return c.Status(500).SendString(errGet.Error())
	}

	// Отправляем ответы клиенту
	return c.JSON(responseData)
}

// SendToKafka Отправка запроса в Kafka
func SendToKafka(data *RequestData) error {
	// Define topic
	topics := map[string]string{
		"users":        "users_requests",
		"accounts":     "accounts_requests",
		"transactions": "transactions_requests",
	}
	topic := topics[data.Service]

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokers := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Fail to start Sarama producer:", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonData),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return nil
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

// GetFromKafka Getting data from topic
func GetFromKafka(data *RequestData, response *ResponseData) error {
	topics := map[string]string{
		"users":        "users_responses",
		"accounts":     "accounts_responses",
		"transactions": "transactions_responses",
	}
	topic := topics[data.Service]
	// Consumer config
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Return.Errors = true
	// Broker config
	brokers := []string{"localhost:9092"}
	// Create a new consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalln("Fail to start Sarama consumer:", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("fail to start Sarama partition consumer: %w", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Message is received from topic(%s)/partition(%d)/offset(%d)\n", topic, msg.Partition, msg.Offset)
			err = json.Unmarshal(msg.Value, response)
			if err != nil {
				return fmt.Errorf("fail to unmarshal message: %w", err)
			}
			log.Printf("Response data: %+v\n", *response) // Проверка данных, полученных из Kafka
		case <-time.After(10 * time.Second): // Периодичность получения ответа
			log.Println("Timeout")
			return nil
		}
	}
}
