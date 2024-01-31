package gateway

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func main() {
	// Создаем новый экземпляр Fiber
	app := fiber.New()

	// Middleware for check auth token
	app.Use(func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" {
			return c.Redirect("/login")
		}

		return c.Next()
	})

	// Определяем обработчик для эндпойнта /api - для примера
	app.Get("/api", func(c *fiber.Ctx) error {
		// Получаем URL удаленного сервера, к которому будем проксировать запрос
		remoteURL := "http://localhost:8080"

		// Создаем новый запрос на основе текущего запроса от клиента
		req := c.Request()

		// Create new and copy the headers manually
		var remoteReqHeader fasthttp.RequestHeader
		req.Header.VisitAll(func(key, value []byte) {
			remoteReqHeader.AddBytesKV(key, value)
		})

		// Use fasthttp to perform the request to the remote server
		statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s%s", remoteURL, req.RequestURI()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error making request to remote server")
		}

		// Copy the status code and response body from the remote server to client
		c.Status(statusCode)
		err1 := c.SendString("Proxy Response:\n")
		if err1 != nil {
			return err1
		}
		err2 := c.SendString(fmt.Sprintf("Status Code: %d\n", statusCode))
		if err2 != nil {
			return err2
		}
		err3 := c.SendString("Body:\n")
		if err3 != nil {
			return err3
		}
		err4 := c.Send(body)
		if err4 != nil {
			return err4
		}

		return nil
	})

	go func() {
		if err := fasthttp.ListenAndServe(":8080", app.Handler()); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()
	select {}
}
