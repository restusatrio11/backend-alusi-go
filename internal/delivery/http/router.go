package http

import (
	"backend-alusi-go/config"
	"backend-alusi-go/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and configures the Gin engine
func SetupRouter(cfg *config.Config, healthHandler *HealthHandler) *gin.Engine {
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Base Middlewares
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Health check route
	router.GET("/healthz", healthHandler.Check)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	return router
}
