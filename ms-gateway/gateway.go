package gateway

import (
	"encoding/json"
	"fmt"
	"log"

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
	data := RequestData{
		Service: "users",
		Action:  "getUsers",
		Data:    map[string]interface{}{"id": "1"},
	}
	log.Printf("Request: %v\n", data)

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Отправляем сообщения в топик
	errSend := SendToKafka(&data)
	if errSend != nil {
		return c.Status(500).SendString(errSend.Error())
	}

	err := SomeWork(&data)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(data)
}

// Мок функция
func SomeWork(data *RequestData) error {
	runes := []rune(data.Action)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	data.Action = string(runes)
	return nil
}

// Отправка запроса в Kafka
func SendToKafka(data *RequestData) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	// config.Producer.Return.Errors = true
	// config.Consumer.Return.Errors = true

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
