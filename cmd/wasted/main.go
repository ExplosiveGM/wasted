package main

import (
	"fmt"
	"log"

	"github.com/ExplosiveGM/wasted/config"
	"github.com/ExplosiveGM/wasted/internal/api"
	"github.com/ExplosiveGM/wasted/internal/db/client"
	"github.com/ExplosiveGM/wasted/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	cfg, _ := config.Load()
	fmt.Println(cfg)
	db, err := client.Connect(&cfg.Database)

	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	defer db.Close()

	logger := logger.NewLogger(cfg)

	router := api.Router(db, logger, cfg)
	logger.Info().Str("env", cfg.App.Env).Msg("🚀 Запуск приложения")

	if err := router.Run(); err != nil {
		logger.Panic().Msg(
			fmt.Sprintf("failed to run server: %v", err),
		)
	}
}
