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

var rootCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management CLI",
	Long:  `Manage database: create, drop, migrate, seed, etc.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.Load()
	},
}

var dbCreateCmd = &cobra.Command{
	Use:   "db-create",
	Short: "Create database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := admin.CreateDatabase(); err != nil {
			log.Fatal("Failed to create database. ", err)
		}
	},
}

var dbDropCmd = &cobra.Command{
	Use:   "db-drop",
	Short: "Drop database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := admin.DropDatabase(); err != nil {
			log.Fatal("Failed to drop database. ", err)
		}
	},
}

var migrationCreateCmd = &cobra.Command{
	Use:   "migration-create",
	Short: "Create migration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		if err := migrate.Create(args[0]); err != nil {
			log.Fatal("Failed to create migration. ", err)
		}

		log.Println("✓ Migration files created")
	},
}

var migrationUpCmd = &cobra.Command{
	Use:   "migration-up",
	Short: "Apply all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrate.Up(); err != nil {
			log.Fatal("Failed apply all migrations. ", err)
		}
	},
}

var migrationDownCmd = &cobra.Command{
	Use:   "migration-down",
	Short: "Rollback to previous migration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrate.Down(); err != nil {
			log.Fatal("Failed rollback to previous migration. ", err)
		}
	},
}

var migrationStatusCmd = &cobra.Command{
	Use:   "migration-status",
	Short: "Print the status of all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrate.Status(); err != nil {
			log.Fatal("Failed print the status of all migrations. ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCreateCmd, dbDropCmd, migrationCreateCmd, migrationUpCmd, migrationDownCmd, migrationStatusCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
