package admin

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func CreateDatabase() error {
	adminDB, err := Connect()
	if err != nil {
		return fmt.Errorf("connect to admin db: %w", err)
	}
	defer adminDB.Close()

	dbName := viper.GetString("DB_NAME")

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`

	err = adminDB.QueryRow(query, dbName).Scan(&exists)

	if err != nil {
		return fmt.Errorf("check database existence: %w", err)
	}

	if !exists {
		_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))

		if err != nil {
			return fmt.Errorf("create database %w", err)
		}

		log.Printf("✓ Database %s created", dbName)
	} else {
		log.Printf("Database %s already exists", dbName)
	}

	return nil
}

func DropDatabase() error {
	adminDB, err := Connect()
	if err != nil {
		return fmt.Errorf("connect to admin db: %w", err)
	}
	defer adminDB.Close()

	dbName := viper.GetString("DB_NAME")

	_, err = adminDB.Exec(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = $1
		AND pid <> pg_backend_pid()`,
		dbName)

	if err != nil {
		log.Printf("Warning: failed to terminate connections: %v", err)
	}

	_, err = adminDB.Exec(fmt.Sprintf("DROP DATABASE %s", dbName))
	if err != nil {
		return fmt.Errorf("drop database: %w", err)
	} else {
		log.Printf("✓ Database %s dropped", dbName)
	}

	return nil
}
