package http

import (
	"backend-alusi-go/config"
	"backend-alusi-go/internal/delivery/http/middleware"
	"backend-alusi-go/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and configures the Gin engine
func SetupRouter(
	cfg *config.Config,
	healthHandler *HealthHandler,
	authHandler *AuthHandler,
	jwtService *jwt.JWTService,
) *gin.Engine {
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
		// Auth Routes (SSO Sumut)
		auth := v1.Group("/auth")
		{
			auth.GET("/login", authHandler.Login)
			auth.GET("/callback", authHandler.Callback)
			auth.GET("/logout", authHandler.Logout)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthRequired(jwtService), authHandler.Me)
		}
	}

	return router
}
