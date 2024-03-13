package main

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func UserAction() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: возвращаем этот код обратно, будет принимать параметры из тела внешнего запроса к шлюзу
		err := proxy.Balancer(proxy.Config{
			Servers: []string{
				// Получаем URL удаленного сервера, к которому будем проксировать запрос
				"http://localhost:9090",
			},
			// Тут какая разница - мы вызываем конкретный путь users - зачем резать путь
			ModifyRequest: func(c *fiber.Ctx) error {
				c.Request().URI().SetPath(strings.TrimPrefix(string(c.Request().URI().Path()), "/user"))
				return nil
			},
			ModifyResponse: func(c *fiber.Ctx) error {
				var responseData string
				if c.Response().StatusCode() >= 400 {
					responseData = "User Service: Error - " + string(c.Response().Header.StatusMessage())
				} else if c.Response().Body() != nil && len(c.Response().Body()) > 0 {
					responseData = "User Service - Users: " + string(c.Response().Body())
				} else {
					responseData = "User Service: No data"
				}
				return c.SendString(responseData)
			},
		})(c)

		// If error with remote service
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("User Service Error - " + err.Error())
		}

		return nil
	}
}
