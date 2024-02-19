package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bsanzhiev/bahamas/user/controllers"
	"github.com/bsanzhiev/bahamas/user/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

	// Группировка роутов
	userController := &controllers.UserController{
		DBPool: dbPool,
		Ctx:    ctx,
	}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/users", func(c *fiber.Ctx) error {
		return c.SendString("Return all accounts v1")
	})
	controllers.RegisterUserRoutes(v1, userController)

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
