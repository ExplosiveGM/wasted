package migrate

import (
	"fmt"

	"github.com/ExplosiveGM/wasted/config"
	"github.com/ExplosiveGM/wasted/internal/db/client"
	"github.com/pressly/goose/v3"
)

func Create(name string) error {
	db, err := client.Connect()

	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer db.Close()

	return goose.Create(db, config.DBMigrations, name, "sql")
}

func Up() error {
	db, err := client.Connect()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, config.DBMigrations); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func Down() error {
	db, err := client.Connect()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer db.Close()

	if err := goose.Down(db, config.DBMigrations); err != nil {
		return fmt.Errorf("rollback migration: %w", err)
	}

	return nil
}

func Status() error {
	db, err := client.Connect()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer db.Close()

	return goose.Status(db, config.DBMigrations)
}
