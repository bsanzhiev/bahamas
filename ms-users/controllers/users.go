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

type IncomingData struct {
	Service string                 `json:"service"` // Имя сервиса
	Action  string                 `json:"action"`  // Имя операции
	Data    map[string]interface{} `json:"data"`    // Объект запроса
}

type OutgoingData struct {
	Status  int         `json:"status"`  // Код ответа
	Message string      `json:"message"` // Сообщение об ошибке или результате
	Data    interface{} `json:"data"`    // Объект ответа
}

// GetUsers How Kafka push to work this code?
func (uc *UserController) GetUsers() ([]User, error) {
	ctx := uc.Ctx
	dbPool := uc.DBPool
	// Делаем запрос
	rows, errRows := dbPool.Query(ctx, "SELECT id, username, first_name, last_name, email FROM users")
	if errRows != nil {
		return nil, errRows
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if errUsers := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email); errUsers != nil {
			return nil, errUsers
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *UserController) getUser() (User, error) {
	ctx := uc.Ctx
	dbPool := uc.DBPool

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
