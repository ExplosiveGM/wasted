package main

import (
	"fmt"
	"log"

	"github.com/ExplosiveGM/wasted/config"
	"github.com/ExplosiveGM/wasted/internal/db/admin"
	"github.com/ExplosiveGM/wasted/internal/db/migrate"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
)

func main() {
	cfg, _ := config.Load()
	rootCmd := NewRootCmd()
	dbCreateCmd := NewDbCreateCmd(&cfg.Database)
	dbDropCmd := NewDbDropCmd(&cfg.Database)
	migrationCreateCmd := NewMigrationCreateCmd(&cfg.Database)
	migrationUpCmd := NewMigrationUpCmd(cfg)
	migrationDownCmd := NewMigrationDownCmd(cfg)
	migrationStatusCmd := NewMigrationStatucCmd(cfg)

	rootCmd.AddCommand(dbCreateCmd, dbDropCmd, migrationCreateCmd, migrationUpCmd, migrationDownCmd, migrationStatusCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "db",
		Short: "Database management CLI",
		Long:  `Manage database: create, drop, migrate, seed, etc.`,
	}
}

func NewDbCreateCmd(dbConfig *config.DatabaseConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "db-create",
		Short: "Create database",
		Run: func(cmd *cobra.Command, args []string) {
			if err := admin.CreateDatabase(dbConfig); err != nil {
				log.Fatal("Failed to create database. ", err)
			}
		},
	}
}

func NewDbDropCmd(dbConfig *config.DatabaseConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "db-drop",
		Short: "Drop database",
		Run: func(cmd *cobra.Command, args []string) {
			if err := admin.DropDatabase(dbConfig); err != nil {
				log.Fatal("Failed to drop database. ", err)
			}
		},
	}
}

func NewMigrationCreateCmd(dbConfig *config.DatabaseConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "migration-create",
		Short: "Create migration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
			if err := migrate.Create(args[0], dbConfig); err != nil {
				log.Fatal("Failed to create migration. ", err)
			}

			log.Println("✓ Migration files created")
		},
	}
}

func NewMigrationUpCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migration-up",
		Short: "Apply all migrations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := migrate.Up(cfg); err != nil {
				log.Fatal("Failed apply all migrations. ", err)
			}
		},
	}
}

func NewMigrationDownCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migration-down",
		Short: "Rollback to previous migration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := migrate.Down(cfg); err != nil {
				log.Fatal("Failed rollback to previous migration. ", err)
			}
		},
	}
}

func NewMigrationStatucCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migration-status",
		Short: "Print the status of all migrations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := migrate.Status(cfg); err != nil {
				log.Fatal("Failed print the status of all migrations. ", err)
			}
		},
	}
}
