package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Alive(c *fiber.Ctx) error {
	defer func() {
		c.JSON(fiber.Map{"alive": true, "ready": true})
	}()
	return nil
}

func main() {
	app := fiber.New()

	app.Get("/", Alive)

	if err := app.Listen(":9091"); err != nil {
		fmt.Printf("Error starting User server: %s\n", err)
	}
}
