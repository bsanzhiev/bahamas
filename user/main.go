package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bsanzhiev/bahamas/user/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
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

func GetUserData(c *fiber.Ctx, ctx context.Context, dbPool *pgxpool.Pool) error {

	// Выполняем запрос
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

func main() {
	ctx := context.Background()
	urlDB := "postgres://postgres:pass123@localhost:9010/bahamas_users"

	app := fiber.New()

	// Create Pool
	dbPool, errPool := pgxpool.New(ctx, urlDB)
	if errPool != nil {
		log.Fatalf("Failed to create pool: %v", errPool)
	}
	defer dbPool.Close()

	fmt.Println("Successfully connected to database!")

	// Делаем миграцию
	migrator, err := migrations.NewMigrator(ctx, dbPool)
	if err != nil {
		panic(err)
	}

	// get the current migration status
	now, exp, info, err := migrator.Info()
	if err != nil {
		panic(err)
	}
	if now < exp {
		// migration is required, dump out the current state
		// and perform the migration
		println("migration needed, current state:")
		println(info)

		err = migrator.Migrate(ctx)
		if err != nil {
			panic(err)
		}
		println("migration successful!")
	} else {
		println("no database migration needed")
	}

	app.Get("/alive", Alive)

	app.Get("/", func(c *fiber.Ctx) error {
		return GetUserData(c, ctx, dbPool)
	})
	if err := app.Listen(":9090"); err != nil {
		fmt.Printf("Error starting User server: %s\n", err)
	}
}
