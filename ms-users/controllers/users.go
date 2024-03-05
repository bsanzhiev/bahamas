package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type UserController struct {
	DBPool *pgxpool.Pool
	Ctx    context.Context
}

func RegisterUserRoutes(app fiber.Router, uc *UserController) {
	users := app.Group("/users")
	users.Get("/list", uc.GetUsers)
	users.Get("/get/:id", getUser)
	users.Post("/create", createUser)
	users.Put("/update/:id", updateUser)
	users.Delete("/delete/:id", deleteUser)
}

func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	// Определяем дополнительно контекст и пул соединений для читаемости
	dbPool := uc.DBPool
	ctx := uc.Ctx

	// Делаем запрос
	rows, errRows := dbPool.Query(ctx, "SELECT id, username, first_name, last_name, email FROM users")
	if errRows != nil {
		return errRows
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if errUsers := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email); errUsers != nil {
			return errUsers
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// получить список пользователей и отправить в ответе
	// return c.SendString("User Data\n")
	return c.JSON(users)
}

func getUser(c *fiber.Ctx) error {
	// get logic here
	return c.SendString("Get user by id")
}

func createUser(c *fiber.Ctx) error {
	// create logic here
	return c.SendString("Create user")
}

func updateUser(c *fiber.Ctx) error {
	// update logic here
	return c.SendString("Update user")
}

func deleteUser(c *fiber.Ctx) error {
	// delete logic here
	return c.SendString("Delete user")
}
