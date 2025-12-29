package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	Jwt      JWTConfig      `mapstructure:",squash"`
	Log      LogConfig      `mapstructure:",squash"`
	Path     PathConfig
}

type AppConfig struct {
	Name string `mapstructure:"APP_NAME"`
	Env  string `mapstructure:"APP_ENV"`
	Port string `mapstructure:"APP_PORT"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SslMode  string `mapstructure:"DB_SSL_MODE"`
}

type JWTConfig struct {
	AccessSecret  string `mapstructure:"JWT_ACCESS_SECRET"`
	RefreshSecret string `mapstructure:"JWT_REFRESH_SECRET"`
}

type LogConfig struct {
	Level       string `mapstructure:"LOG_LEVEL"`
	File        string `mapstructure:"LOG_FILE"`
	EnableJson  bool   `mapstructure:"LOG_ENABLE_JSON"`
	EnableColor bool   `mapstructure:"LOG_ENABLE_COLOR"`
}

type PathConfig struct {
	RootDir      string
	DBDir        string
	DBMigrations string
}

var (
	RootDir      string
	DBDir        string
	DBMigrations string
)

func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")

	if strings.ToLower(env) == "test" {
		viper.SetConfigName(".env.test")
	} else {
		viper.SetConfigName(".env")
	}

	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println(".env file not found, using environment variables")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	} else if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	setPaths(&cfg.Path)

	return &cfg, nil
}

func setPaths(pathConfig *PathConfig) {
	rootDir := findProjectRootByGoMod()
	dbDir := filepath.Join(rootDir, "db")
	dbMigrations := filepath.Join(dbDir, "migrations")

	pathConfig.RootDir = rootDir
	pathConfig.DBDir = dbDir
	pathConfig.DBMigrations = dbMigrations
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
