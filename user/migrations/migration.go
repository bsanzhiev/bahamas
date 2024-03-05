package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

const versionTable = "schema_migrations"

// Используем встроенную файловую структуру для хранения скриптов SQL

//go:embed scripts/*
var migrationsFS embed.FS

type Migrator struct {
	migrator *migrate.Migrator
}

func NewMigrator(ctx context.Context, dbPool *pgxpool.Pool) (*Migrator, error) {
	// Берем одно соединение
	conn, errConn := dbPool.Acquire(ctx)
	if errConn != nil {
		return nil, fmt.Errorf("failed to acquire a connection from the pool: %v", errConn)
	}
	defer conn.Release()

	// Создаем новый экземпляр
	migrator, errNewMigrator := migrate.NewMigrator(ctx, conn.Conn(), versionTable)
	if errNewMigrator != nil {
		return nil, fmt.Errorf("failed to create migrator: %v", errNewMigrator)
	}

	// Загружаем миграции из файлов
	scripts, errScripts := fs.Sub(migrationsFS, "scripts")
	if errScripts != nil {
		return nil, fmt.Errorf("failed to load migrations: %v", errScripts)
	}
	// scripts, _ := fs.Sub(migrationsFS, "scripts")

	// Делаем миграцию в базу данных
	if errMigration := migrator.LoadMigrations(scripts); errMigration != nil {
		return nil, fmt.Errorf("failed to load migrations: %v", errMigration)
	}

	return &Migrator{migrator: migrator}, nil
}

// Миграция
func (m *Migrator) Migrate(ctx context.Context) error {
	if err := m.migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}
	fmt.Println("Migrations applied successfully")
	return nil
}

// Info the current migration version and the embedded maximum migration, and a textual
// representation of the migration state for informational purposes.
func (m Migrator) Info() (int32, int32, string, error) {

	version, err := m.migrator.GetCurrentVersion(context.Background())
	if err != nil {
		return 0, 0, "", err
	}
	info := ""

	var last int32
	for _, thisMigration := range m.migrator.Migrations {
		last = thisMigration.Sequence

		cur := version == thisMigration.Sequence
		indicator := "  "
		if cur {
			indicator = "->"
		}
		info = info + fmt.Sprintf(
			"%2s %3d %s\n",
			indicator,
			thisMigration.Sequence, thisMigration.Name)
	}

	return version, last, info, nil
}

// MigrateTo migrates to a specific version of the schema. Use '0' to
// undo all migrations.
func (m Migrator) MigrateTo(ver int32) error {
	err := m.migrator.MigrateTo(context.Background(), ver)
	if err != nil {
		return fmt.Errorf("migration to version %d failed: %v", ver, err)
	}
	return nil // Возвращаем nil в случае успеха
}
