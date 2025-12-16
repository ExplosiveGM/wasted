package main

import (
	"fmt"
	"log"

	"github.com/ExplosiveGM/wasted/config"
	"github.com/ExplosiveGM/wasted/internal/api"
	"github.com/ExplosiveGM/wasted/internal/db/client"
	"github.com/ExplosiveGM/wasted/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

func main() {
	config.Load()

	db, err := client.Connect()

	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	defer db.Close()

	logger := logger.NewLogger(logger.Config{
		Environment: viper.GetString("APP_ENV"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
		LogFile:     viper.GetString("LOG_FILE"),
		EnableJSON:  viper.GetBool("LOG_ENABLE_JSON"),
		EnableColor: viper.GetBool("LOG_ENABLE_COLOR"),
	})

	router := api.Router(db, logger)

	logger.Info().Str("env", viper.GetString("APP_ENV")).Msg("🚀 Запуск приложения")

	if err := router.Run(); err != nil {
		logger.Panic().Msg(
			fmt.Sprintf("failed to run server: %v", err),
		)
	}
}
