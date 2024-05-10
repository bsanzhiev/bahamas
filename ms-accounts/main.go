package main

import (
	"context"
	"fmt"
	"github.com/bsanzhiev/bahamas/ms-accounts/migrations"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Database connection ========================
	// Get connection string from env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot get env file.")
		return
	}
	connString := os.Getenv("CONNECTION_STRING")

	// Create connections pool
	urlDB := connString
	dbPool, err := pgxpool.New(ctx, urlDB)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer dbPool.Close()
	fmt.Println("Successfully connected to Accounts database!")

	// Migration =======================================
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

	// Logic for work with database and Kafka
	// ...

	// Define Fiber app =======================
	app := fiber.New(
		fiber.Config{
			AppName: "Bahamas Accounts Service",
		})
	// Check alive available
	app.Get("/alive", Alive)
	// Start app
	if err := app.Listen(":7003"); err != nil {
		fmt.Printf("Error starting Account server: %s\n", err)
	}
}

func Alive(c *fiber.Ctx) error {
	defer func() {
		err := c.JSON(fiber.Map{"alive": true, "ready": true, "service": "accounts"})
		if err != nil {
			return
		}
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
