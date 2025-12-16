package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	RootDir      string
	DBDir        string
	DBMigrations string
)

func Load() {
	RootDir = findProjectRootByGoMod()
	DBDir = filepath.Join(RootDir, "db")
	DBMigrations = filepath.Join(DBDir, "migrations")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println(".env file not found, using environment variables")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	viper.AutomaticEnv()
}

func findProjectRootByGoMod() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		goModPath := filepath.Join(cwd, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return cwd
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			break // Reached filesystem root
		}
		cwd = parent
	}

	return ""
}

func findEnvFile() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
}
