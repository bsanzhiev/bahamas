package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func TransactionController(app fiber.Router) {
	accounts := app.Group("/transactions")
	accounts.Get("/list", getTransactions)
	accounts.Get("/get/:id", getTransaction)
	accounts.Post("/create", createTransaction)
	accounts.Delete("/delete/:id", deleteTransaction)
}

func getTransactions(c *fiber.Ctx) error {
	// get logic here
	return c.SendString("List of transaction")
}

func getTransaction(c *fiber.Ctx) error {
	// get logic here
	return c.SendString("Get transaction by id")
}

func createTransaction(c *fiber.Ctx) error {
	// create logic here
	return c.SendString("Create transaction")
}

// Update is not defined

func deleteTransaction(c *fiber.Ctx) error {
	// delete logic here
	return c.SendString("Delete trancaction")
}
