package gateway

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func StartGateway() {
	// Создаем новый экземпляр Fiber
	app := fiber.New()

	// Middleware for check auth token
	// app.Use(func(c *fiber.Ctx) error {
	// 	token := c.Get("Authorization")

	// 	if token == "" {
	// 		c.Status(401).JSON(fiber.Map{"error": "No valid token"})
	// 		return c.Redirect("/login")
	// 	}
	// 	return c.Next()
	// })

	app.Get("/", Alive)

	// User service
	app.Use("/users/*", UserAction())

	// Запуск сервера шлюза
	go func() {
		if err := app.Listen(":9080"); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()
	select {}
}

func Alive(c *fiber.Ctx) error {
	defer func() {
		c.JSON(fiber.Map{"alive": true, "ready": true, "service": "gateway"})
	}()
	return nil
}
