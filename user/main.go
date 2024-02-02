package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	postgres "github.com/gofiber/storage/postgres/v2"
)

func Alive(c *fiber.Ctx) error {
	defer func() {
		c.JSON(fiber.Map{"alive": true, "ready": true})
	}()
	return nil
}

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

func GetUserData(c *fiber.Ctx) error {

	ConfigDefault := postgres.Config{
		ConnectionURI: "",
		Host:          "127.0.0.1",
		Port:          9010,
		Username:      "postgres",
		Password:      "pass123",
		Database:      "bahamas_users",
		Table:         "users",
		SSLMode:       "disable",
		Reset:         false,
		GCInterval:    10 * time.Second,
	}
	store := postgres.New(ConfigDefault)

	db, err := sql.Open("postgres", store.Conn().Config().ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	// Выполняем запрос
	rows, errRows := db.Query("SELECT * FROM users")
	if errRows != nil {
		return errRows
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if errUsers := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.Username, &user.FirstName, &user.LastName, &user.Email); errUsers != nil {
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

func main() {
	app := fiber.New()

	app.Get("/alive", Alive)

	app.Get("/", GetUserData)

	if err := app.Listen(":9090"); err != nil {
		fmt.Printf("Error starting User server: %s\n", err)
	}
}
