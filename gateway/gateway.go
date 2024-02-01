package gateway

// TODO - переписать без использования fasthttp

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func StartGateway() {
	// Создаем новый экземпляр Fiber
	app := fiber.New()

	// Middleware for check auth token
	//app.Use(func(c *fiber.Ctx) error {
	//	token := c.Get("Authorization")
	//
	//	if token == "" {
	//		return c.Redirect("/login")
	//	}
	//
	//	return c.Next()
	//})

	// User - определяем обработчик для эндпойнта /user - для примера
	app.Use("/user/*", proxy.Balancer(proxy.Config{
		Servers: []string{
			// Получаем URL удаленного сервера, к которому будем проксировать запрос
			"http://localhost:9090",
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			// Есть ли тело ответа от удаленного сервиса
			// Ставим заглушку
			//if c.Response().Body() != nil && len(c.Response().Body()) > 0 {
			//	responseData := string(c.Response().Body()) + "User Data"
			//	return c.SendString(responseData)
			//} else {
			return c.SendString("User data")
			//}
		},
	}))

	// Account
	app.Use("/account", proxy.Balancer(proxy.Config{
		Servers: []string{
			// Получаем URL удаленного сервера, к которому будем проксировать запрос
			"http://localhost:9091",
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			return nil
		},
	}))

	// Transactions
	app.Use("/transaction", proxy.Balancer(proxy.Config{
		Servers: []string{
			// Получаем URL удаленного сервера, к которому будем проксировать запрос
			"http://localhost:9092",
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			return nil
		},
	}))

	go func() {
		if err := app.Listen(":9080"); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()
	select {}
}
