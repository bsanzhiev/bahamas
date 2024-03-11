package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
)

func StartGateway() {
	// Создаем новый экземпляр Fiber
	app := fiber.New()

	// Сервис запущен
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

	// Что приходит в запросе:
	// - Имя сервиса
	// - Имя операции
	// - Объект запроса

	// Разбираем полученный запрос
	// Узнаем имя сервиса - сопоставляем с адресами конкретного сервиса
	// Отправляем данные запроса дальше

	// имена сервисов лежат в карте - сопостовляем с топиками
	// Используем Кафку
	// Кидаем в нужный топик пришедшие данные
	// Получаем данные из топика
	// Отправляем обратно клиенту ответ
	// Как получить данные обратно?
	// Для отправки и получения сообщений используем разные топики

	app.Post("/", HandleRequest)

	// User service
	// app.Use("/users/*", UserAction())

	// Запуск шлюза
	go func() {
		if err := app.Listen(":9080"); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()
	select {}
}

// Проверка работы
func Alive(c *fiber.Ctx) error {
	defer func() {
		c.JSON(fiber.Map{"alive": true, "ready": true, "service": "gateway"})
	}()
	return nil
}

// Тип для запросов
type RequestData struct {
	Service string                 `json:"service"`
	Action  string                 `json:"action"`
	Data    map[string]interface{} `json:"data"`
}

// Обработка входящих запросов
func HandleRequest(c *fiber.Ctx) error {
	var data RequestData

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	err := SomeWork(&data)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(data)
}

func SomeWork(data *RequestData) error {
	runes := []rune(data.Action)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	data.Action = string(runes)
	return nil
}

func SendToKafka(data *RequestData) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_1_0_0

	brokers := []string{"localhost:9100"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Fail to start Sarama producer:", err)
		os.Exit(1)
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
		Topic: "users_requests",
		Value: sarama.ByteEncoder(jsonData),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return nil
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "users_requests", partition, offset)
	return nil
}
