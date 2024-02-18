package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bsanzhiev/bahamas/account/controllers"
	"github.com/bsanzhiev/bahamas/account/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	app.Get("/alive", Alive)

	// Создаем пул подключений к базе данных =============================
	ctx := context.Background()
	urlDB := "postgres://postgres:pass123@localhost:9011/bahamas_accounts"

	dbPool, errPool := pgxpool.New(ctx, urlDB)
	if errPool != nil {
		log.Fatalf("Failed to create pool: %v", errPool)
	}
	defer dbPool.Close()
	fmt.Println("Successfully connected to Accounts database!")

	// Делаем миграцию ========================================
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
	// ======================================

	// Здесь будет группа роутов
	// 1 - Счета - группа - список счетов, создание счета, изменение счета, удаление
	// 2 - Транзакции - группа - тоже самое
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/accounts", func(c *fiber.Ctx) error {
		return c.SendString("Return all accounts v1")
	})
	controllers.AccountController(v1)
	v1.Get("/transactions", func(c *fiber.Ctx) error {
		return c.SendString("Return all transactions v1")
	})
	controllers.TransactionController(v1)

	// Старт сервиса
	if err := app.Listen(":9091"); err != nil {
		fmt.Printf("Error starting User server: %s\n", err)
	}
}

func Alive(c *fiber.Ctx) error {
	defer func() {
		c.JSON(fiber.Map{"alive": true, "ready": true})
	}()
	return nil
}

type Account struct {
	AccountID     uuid.UUID `json:"account_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	DeletedAt     time.Time `json:"deleted_at,omitempty"`
	UserID        int       `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	AccountType   string    `json:"account_type"`
	Balance       float64   `json:"balance"`
}

type Transaction struct {
	TransactionID      uuid.UUID  `json:"transaction_id"`       // Идентификатор транзакции в формате UUID
	CreatedAt          time.Time  `json:"created_at"`           // Время создания транзакции
	UpdatedAt          *time.Time `json:"updated_at,omitempty"` // Время последнего обновления транзакции, может быть nil
	DeletedAt          *time.Time `json:"deleted_at,omitempty"` // Время удаления транзакции, может быть nil
	AccountID          uuid.UUID  `json:"account_id"`           // ID счета отправителя
	RecipientAccountID int        `json:"recipient_account_id"` // ID счета получателя
	Amount             float64    `json:"amount"`               // Сумма транзакции
	TransactionType    string     `json:"transaction_type"`     // Тип транзакции (например, "пополнение", "снятие")
	TransactionDate    time.Time  `json:"transaction_date"`     // Дата и время проведения транзакции
	Status             string     `json:"status"`               // Статус транзакции
	Description        string     `json:"description"`          // Описание транзакции
}

type AccountHolder struct {
	AccountID uuid.UUID `json:"account_id"` // Внешний ключ, связывающий с таблицей "accounts"
	UserID    int       `json:"user_id"`    // Внешний ключ, связывающий с таблицей "users"
}
