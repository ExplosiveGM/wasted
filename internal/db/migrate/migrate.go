package migrate

import (
	"fmt"

	"github.com/ExplosiveGM/wasted/config"
	"github.com/ExplosiveGM/wasted/internal/db/client"
	"github.com/pressly/goose/v3"
)

func Create(name string, dbConfig *config.DatabaseConfig) error {
	db, err := client.Connect(dbConfig)

	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	defer db.Close()

	return goose.Create(db, config.DBMigrations, name, "sql")
}

func Up(cfg *config.Config) error {
	db, err := client.Connect(&cfg.Database)

	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	defer db.Close()

	if err := goose.Up(db, cfg.Path.DBMigrations); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func Down(cfg *config.Config) error {
	db, err := client.Connect(&cfg.Database)

	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	defer db.Close()

	if err := goose.Down(db, cfg.Path.DBMigrations); err != nil {
		return fmt.Errorf("rollback migration: %w", err)
	}

	return nil
}

func Status(cfg *config.Config) error {
	db, err := client.Connect(&cfg.Database)

	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	defer db.Close()

	return goose.Status(db, cfg.Path.DBMigrations)
}
