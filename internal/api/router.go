package api

import (
	"database/sql"

	"github.com/ExplosiveGM/wasted/docs"
	"github.com/ExplosiveGM/wasted/internal/auth"
	"github.com/ExplosiveGM/wasted/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(db *sql.DB, logger zerolog.Logger) *gin.Engine {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
