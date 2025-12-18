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

//	@title			wasted
//	@version		1.0
//	@description	Wasted API.
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:8080
//	@BasePath		/api/v1
//	@schemes		http https

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
