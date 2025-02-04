package postgres

import (
	"customers_service/internal/infrastructure/config"
	"database/sql"
	"fmt"
	"log"
)

type Repository struct {
	cfg config.PostgresConfig
	db  *sql.DB
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewRepository(cfg config.PostgresConfig) *Repository {
	repo := &Repository{
		cfg: cfg,
	}

	repo.init()

	return repo
}

func (r *Repository) init() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		r.cfg.Host, r.cfg.Port, r.cfg.User, r.cfg.Password, r.cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	r.db = db
}

func (r *Repository) Close() {
	if err := r.db.Close(); err != nil {
		log.Printf("failed to close postgres connection: %v", err)
	}
}
