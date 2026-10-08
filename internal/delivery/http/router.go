package http

import (
	"fmt"
	"strings"
	"time"

	"backend-alusi-go/config"
	"backend-alusi-go/docs"
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
	rbacHandler *RBACHandler,
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

	// Root & Health check routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"app":     cfg.App.Name,
			"version": "1.2.0",
			"message": "ALUSI Backend API is operational",
			"docs":    "/swagger/index.html",
		})
	})
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"app":     cfg.App.Name,
			"version": "1.2.0",
			"message": "ALUSI Backend API is operational",
		})
	})
	router.GET("/api/index.go", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"app":     cfg.App.Name,
			"version": "1.2.0",
			"message": "ALUSI Serverless Backend API is operational",
		})
	})
	router.GET("/healthz", healthHandler.Check)

	// Interactive Swagger API Documentation
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	// Static file serving for uploads (logos, attachments)
	router.Static("/uploads", "./uploads")

	// Top-level SSO OAuth redirect route (/auth/callback)
	// Forwards code & state from SSO provider redirect (http://localhost:8080/auth/callback) to frontend or auth handler
	router.GET("/auth/callback", func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")

		// Determine frontend origin from CORS allowed origins or default to http://localhost:5173
		frontendOrigin := "http://localhost:5173"
		for _, origin := range cfg.CORS.AllowedOrigins {
			if strings.Contains(origin, "5173") || strings.Contains(origin, "3000") || strings.Contains(origin, "bps.web.id") {
				frontendOrigin = origin
				break
			}
		}

		targetURL := fmt.Sprintf("%s/callback?code=%s&state=%s", frontendOrigin, code, state)
		c.Redirect(302, targetURL)
	})

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
		{
			// Analytics & Reporting
			analytics := admin.Group("/analytics")
			analytics.Use(middleware.RequireAnyPermission("analytics:view", "dashboard:view"))
			{
				analytics.GET("/summary", analyticsHandler.GetDashboardSummary)
				analytics.GET("/top-apps", analyticsHandler.GetTopApps)
				analytics.GET("/trends", analyticsHandler.GetTrends)
				analytics.GET("/disruptions", analyticsHandler.GetDisruptions)
			}

			// Executive Report Exports (CSV/Excel)
			reports := admin.Group("/reports")
			reports.Use(middleware.RequirePermission("reports:export"))
			{
				reports.GET("/analytics/export", reportHandler.ExportAnalyticsCSV)
				reports.GET("/catalog/export", reportHandler.ExportCatalogCSV)
			}

			// Apps & Guides CRUD
			admin.POST("/apps", middleware.RequirePermission("apps:create"), adminHandler.CreateApp)
			admin.PUT("/apps/reorder", middleware.RequirePermission("apps:reorder"), adminHandler.ReorderApps)
			admin.PUT("/apps/:id", middleware.RequirePermission("apps:update"), adminHandler.UpdateApp)
			admin.DELETE("/apps/:id", middleware.RequirePermission("apps:delete"), adminHandler.DeleteApp)
			admin.POST("/apps/:id/logo", middleware.RequirePermission("apps:upload_logo"), adminHandler.UploadAppLogo)
			admin.POST("/apps/:id/probe", middleware.RequirePermission("monitoring:probe"), monitoringHandler.ManualProbeApp)

			// App Guides & FAQ Management
			admin.POST("/apps/:id/guides", middleware.RequirePermission("guides:manage"), adminHandler.CreateGuide)
			admin.PUT("/guides/:id", middleware.RequirePermission("guides:manage"), adminHandler.UpdateGuide)
			admin.DELETE("/guides/:id", middleware.RequirePermission("guides:manage"), adminHandler.DeleteGuide)

			// Announcements Broadcast Management
			admin.GET("/announcements", middleware.RequirePermission("announcements:view"), announcementHandler.ListAdminAnnouncements)
			admin.POST("/announcements", middleware.RequirePermission("announcements:create"), announcementHandler.CreateAnnouncement)
			admin.PUT("/announcements/:id", middleware.RequirePermission("announcements:update"), announcementHandler.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", middleware.RequirePermission("announcements:delete"), announcementHandler.DeleteAnnouncement)

			// Feedback & Issue Reports Tracking
			admin.GET("/feedbacks", middleware.RequirePermission("feedbacks:view"), feedbackHandler.ListAdminFeedbacks)
			admin.PUT("/feedbacks/:id", middleware.RequirePermission("feedbacks:respond"), feedbackHandler.UpdateFeedbackStatus)

			// Audit Logs Trail
			admin.GET("/audit-logs", middleware.RequirePermission("audit:view"), auditHandler.ListAuditLogs)

			// Categories Management
			admin.POST("/categories", middleware.RequirePermission("categories:create"), adminHandler.CreateCategory)
			admin.PUT("/categories/:id", middleware.RequirePermission("categories:update"), adminHandler.UpdateCategory)
			admin.DELETE("/categories/:id", middleware.RequirePermission("categories:delete"), adminHandler.DeleteCategory)

			// Dynamic RBAC & Role Permission Management
			rbac := admin.Group("/rbac")
			rbac.Use(middleware.RequireAnyPermission("rbac:view", "users:view"))
			{
				rbac.GET("/permissions", rbacHandler.ListPermissions)
				rbac.GET("/roles", rbacHandler.ListRoles)
				rbac.PUT("/roles/:id/permissions", middleware.RequireAnyPermission("rbac:manage", "rbac:manage_roles"), rbacHandler.UpdateRolePermissions)
				rbac.GET("/users", rbacHandler.ListUsers)
				rbac.PUT("/users/:id/roles", middleware.RequireAnyPermission("users:manage", "rbac:assign_users"), rbacHandler.AssignUserRoles)
			}
		}
	}

	return router
}
