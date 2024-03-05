package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bsanzhiev/bahamas/ms-users/controllers"
	"github.com/bsanzhiev/bahamas/ms-users/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Получаем строку подключения
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}
	connString := os.Getenv("CONNECTION_STRING")

	// urlDB := "postgres://postgres:pass123@localhost:9010/bahamas_users"
	urlDB := connString

	app := fiber.New()

	// Создаем пул подключений к базе данных
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

	// Группировка роутов
	userController := &controllers.UserController{
		DBPool: dbPool,
		Ctx:    ctx,
	}
	api := app.Group("/api/v1")
	// v1 := api.Group("/v1")
	api.Get("/users", func(c *fiber.Ctx) error {
		return c.SendString("Return all accounts v1")
	})
	controllers.RegisterUserRoutes(api, userController)

	if err := app.Listen(":9090"); err != nil {
		fmt.Printf("Error starting User server: %s\n", err)
	}
}

func Alive(c *fiber.Ctx) error {
	defer func() {
		c.JSON(fiber.Map{"alive": true, "ready": true})
	}()
	return nil
}
