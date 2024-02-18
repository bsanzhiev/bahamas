package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func AccountController(app fiber.Router) {
	accounts := app.Group("/accounts")
	accounts.Get("/list", getAccounts)
	accounts.Post("/create", createAccount)
	accounts.Get("/get/:id", getAccount)
	accounts.Put("/update/:id", updateAccount)
	accounts.Delete("/delete/:id", deleteAccount)
}

func getAccounts(c *fiber.Ctx) error {
	// get logic here
	return c.SendString("Get list of accounts")
}

func getAccount(c *fiber.Ctx) error {
	// get logic here
	return c.SendString("Get account by id")
}

func createAccount(c *fiber.Ctx) error {
	// create logic here
	return c.SendString("Create account")
}

func updateAccount(c *fiber.Ctx) error {
	// update logic here
	return c.SendString("Update account")
}

func deleteAccount(c *fiber.Ctx) error {
	// delete logic here
	return c.SendString("Delete account")
}
