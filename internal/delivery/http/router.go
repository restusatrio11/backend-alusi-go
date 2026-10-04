package http

import (
	"time"

	"backend-alusi-go/config"
	"backend-alusi-go/internal/delivery/http/middleware"
	"backend-alusi-go/pkg/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SetupRouter initializes and configures the Gin engine
func SetupRouter(
	cfg *config.Config,
	healthHandler *HealthHandler,
	authHandler *AuthHandler,
	catalogHandler *CatalogHandler,
	interactionHandler *InteractionHandler,
	adminHandler *AdminHandler,
	monitoringHandler *MonitoringHandler,
	announcementHandler *AnnouncementHandler,
	jwtService *jwt.JWTService,
) *gin.Engine {
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Rate Limiter: 100 req/sec with burst 50 per IP
	ipLimiter := middleware.NewIPRateLimiter(rate.Limit(100), 50, 5*time.Minute)

	// Global Base Middlewares
	router.Use(middleware.RequestID())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))
	router.Use(middleware.RateLimit(ipLimiter))

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

		// User Personalization Routes
		users := v1.Group("/users")
		users.Use(middleware.AuthRequired(jwtService))
		{
			users.GET("/me/favorites", interactionHandler.ListFavorites)
			users.GET("/me/recents", interactionHandler.ListRecents)
		}

		// Category Routes
		categories := v1.Group("/categories")
		{
			categories.GET("", catalogHandler.ListCategories)
			categories.GET("/:slug", catalogHandler.GetCategoryBySlug)
		}

		// Application Catalog & Interaction Routes
		apps := v1.Group("/apps")
		{
			apps.GET("", middleware.OptionalAuth(jwtService), catalogHandler.ListApps)
			apps.GET("/search", middleware.OptionalAuth(jwtService), catalogHandler.SearchApps)
			apps.GET("/:slug", middleware.OptionalAuth(jwtService), catalogHandler.GetAppBySlug)
			apps.GET("/:slug/guides", catalogHandler.GetAppGuides)
			apps.GET("/:slug/status-history", monitoringHandler.GetAppStatusHistory)
			apps.POST("/:id/click", middleware.OptionalAuth(jwtService), interactionHandler.RecordClick)
			apps.POST("/:id/favorite", middleware.AuthRequired(jwtService), interactionHandler.ToggleFavorite)
		}

		// Announcements Routes (Public)
		announcements := v1.Group("/announcements")
		{
			announcements.GET("", announcementHandler.GetActiveAnnouncements)
		}

		// Public Service Monitoring
		services := v1.Group("/services")
		{
			services.GET("/status", monitoringHandler.GetServiceUptimeSummary)
		}

		// Admin Management Routes (Protected by Admin Role)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired(jwtService))
		admin.Use(middleware.RequireRoles("admin"))
		{
			admin.POST("/apps", adminHandler.CreateApp)
			admin.PUT("/apps/reorder", adminHandler.ReorderApps)
			admin.PUT("/apps/:id", adminHandler.UpdateApp)
			admin.DELETE("/apps/:id", adminHandler.DeleteApp)
			admin.POST("/apps/:id/probe", monitoringHandler.ManualProbeApp)

			// App Guides & FAQ Management
			admin.POST("/apps/:id/guides", adminHandler.CreateGuide)
			admin.PUT("/guides/:id", adminHandler.UpdateGuide)
			admin.DELETE("/guides/:id", adminHandler.DeleteGuide)

			// Announcements Broadcast Management
			admin.GET("/announcements", announcementHandler.ListAdminAnnouncements)
			admin.POST("/announcements", announcementHandler.CreateAnnouncement)
			admin.PUT("/announcements/:id", announcementHandler.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", announcementHandler.DeleteAnnouncement)

			admin.POST("/categories", adminHandler.CreateCategory)
			admin.PUT("/categories/:id", adminHandler.UpdateCategory)
			admin.DELETE("/categories/:id", adminHandler.DeleteCategory)
		}
	}

	return router
}
