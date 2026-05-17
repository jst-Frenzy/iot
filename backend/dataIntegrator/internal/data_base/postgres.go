package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/config/credentials"
	"github.com/pressly/goose/v3"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(cred *credentials.Credentials) (*gorm.DB, error) {
	postgresDB, err := gorm.Open(
		postgres.Open(cred.PostgresDSN),
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to postgres: %w", err)
	}

	sqlDB, err := postgresDB.DB()
	if err != nil {
		return nil, fmt.Errorf("cannot get sql db: %w", err)
	}

	if err := applyMigrations(sqlDB); err != nil {
		return nil, fmt.Errorf("cannot apply migrations: %w", err)
	}

	return postgresDB, nil
}

func applyMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations/postgres"); err != nil {
		return err
	}

	return nil
}
