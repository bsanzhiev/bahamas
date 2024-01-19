package gateway

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func main() {
	// Создаем новый экземпляр Fiber
	app := fiber.New()

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
		c.SendString("Proxy Response:\n")
		c.SendString(fmt.Sprintf("Status Code: %d\n", statusCode))
		c.SendString("Body:\n")
		c.Send(body)

		return nil
	})

	go func() {
		if err := fasthttp.ListenAndServe(":8080", app.Handler()); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()
	select {}
}
