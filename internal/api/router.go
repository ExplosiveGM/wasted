package api

import (
	"database/sql"

	"github.com/ExplosiveGM/wasted/internal/auth"
	"github.com/ExplosiveGM/wasted/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Router(db *sql.DB, logger zerolog.Logger) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		authRoute := v1.Group("/auth")
		{
			queries := database.New(db)
			authService := auth.NewAuthService(queries, logger)
			authHandler := auth.NewAuthHandler(authService)
			authRoute.POST("/request-code", authHandler.RequestCode)
			authRoute.POST("/verify", authHandler.Verify)
			authRoute.POST("/refresh", authHandler.Refresh)
		}
	}

	return router
}
