package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("User Data\n")
	})

	if err := app.Listen(":9090"); err != nil {
		fmt.Printf("Error starting User server: %s\n", err)
	}
}
