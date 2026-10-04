package http

import (
	"time"

	_ "backend-alusi-go/docs"
	"backend-alusi-go/config"
	"backend-alusi-go/internal/delivery/http/middleware"
	"backend-alusi-go/pkg/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	feedbackHandler *FeedbackHandler,
	analyticsHandler *AnalyticsHandler,
	auditHandler *AuditHandler,
	aiHandler *AIHandler,
	reportHandler *ReportHandler,
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

	// Interactive Swagger API Documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	// Static file serving for uploads (logos, attachments)
	router.Static("/uploads", "./uploads")

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth Routes (SSO Sumut & Manual Login)
		auth := v1.Group("/auth")
		{
			auth.GET("/login", authHandler.Login)
			auth.POST("/login", authHandler.ManualLogin)
			auth.POST("/manual-login", authHandler.ManualLogin)
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

		// Feedbacks & Issue Reports (Public / User)
		feedbacks := v1.Group("/feedbacks")
		{
			feedbacks.POST("", middleware.OptionalAuth(jwtService), feedbackHandler.SubmitFeedback)
		}

		// AI Assistant Recommendation Endpoint (Public / User)
		ai := v1.Group("/ai")
		{
			ai.POST("/ask", middleware.OptionalAuth(jwtService), aiHandler.Ask)
		}

		// Open API Standard Metadata (Public / Satu Data Indonesia)
		openapi := v1.Group("/openapi")
		{
			openapi.GET("/apps", reportHandler.GetOpenAPICatalog)
		}

		// Public Service Monitoring
		services := v1.Group("/services")
		{
			services.GET("/status", monitoringHandler.GetServiceUptimeSummary)
			services.GET("/realtime-status", monitoringHandler.StreamStatusEvents)
		}

		// Admin & Pimpinan Management Routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired(jwtService))
		admin.Use(middleware.RequireRoles("admin", "pimpinan"))
		{
			// Analytics & Reporting (Admin & Pimpinan)
			analytics := admin.Group("/analytics")
			{
				analytics.GET("/summary", analyticsHandler.GetDashboardSummary)
				analytics.GET("/top-apps", analyticsHandler.GetTopApps)
				analytics.GET("/trends", analyticsHandler.GetTrends)
				analytics.GET("/disruptions", analyticsHandler.GetDisruptions)
			}

			// Executive Report Exports (CSV/Excel)
			reports := admin.Group("/reports")
			{
				reports.GET("/analytics/export", reportHandler.ExportAnalyticsCSV)
				reports.GET("/catalog/export", reportHandler.ExportCatalogCSV)
			}

			// Apps, Audit & Category CRUD (Admin Only)
			adminApps := admin.Group("")
			adminApps.Use(middleware.RequireRoles("admin"))
			{
				adminApps.POST("/apps", adminHandler.CreateApp)
				adminApps.PUT("/apps/reorder", adminHandler.ReorderApps)
				adminApps.PUT("/apps/:id", adminHandler.UpdateApp)
				adminApps.DELETE("/apps/:id", adminHandler.DeleteApp)
				adminApps.POST("/apps/:id/logo", adminHandler.UploadAppLogo)
				adminApps.POST("/apps/:id/probe", monitoringHandler.ManualProbeApp)

				// App Guides & FAQ Management
				adminApps.POST("/apps/:id/guides", adminHandler.CreateGuide)
				adminApps.PUT("/guides/:id", adminHandler.UpdateGuide)
				adminApps.DELETE("/guides/:id", adminHandler.DeleteGuide)

				// Announcements Broadcast Management
				adminApps.GET("/announcements", announcementHandler.ListAdminAnnouncements)
				adminApps.POST("/announcements", announcementHandler.CreateAnnouncement)
				adminApps.PUT("/announcements/:id", announcementHandler.UpdateAnnouncement)
				adminApps.DELETE("/announcements/:id", announcementHandler.DeleteAnnouncement)

				// Feedback & Issue Reports Tracking
				adminApps.GET("/feedbacks", feedbackHandler.ListAdminFeedbacks)
				adminApps.PUT("/feedbacks/:id", feedbackHandler.UpdateFeedbackStatus)

				// Audit Logs Trail
				adminApps.GET("/audit-logs", auditHandler.ListAuditLogs)

				adminApps.POST("/categories", adminHandler.CreateCategory)
				adminApps.PUT("/categories/:id", adminHandler.UpdateCategory)
				adminApps.DELETE("/categories/:id", adminHandler.DeleteCategory)
			}
		}
	}

	return router
}
