package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"backend-alusi-go/config"
	deliveryHTTP "backend-alusi-go/internal/delivery/http"
	"backend-alusi-go/internal/repository/postgres"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/ai"
	"backend-alusi-go/pkg/database"
	"backend-alusi-go/pkg/exporter"
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/logger"
	"backend-alusi-go/pkg/media"
	"backend-alusi-go/pkg/realtime"
	"backend-alusi-go/pkg/sso"
	"backend-alusi-go/pkg/worker"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	appRouter *gin.Engine
	initOnce  sync.Once
)

func initApp() {
	// 1. Load Configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load configuration")
		return
	}

	// 2. Initialize Structured Logger
	logger.InitLogger(cfg.App.Debug, cfg.App.Env)

	// 3. Initialize PostgreSQL Database Connection Pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewPostgresDB(ctx, &cfg.Database)

	var userRepo *postgres.UserRepo
	var categoryRepo *postgres.CategoryRepo
	var appRepo *postgres.AppRepo
	var favoriteRepo *postgres.FavoriteRepo
	var clickLogRepo *postgres.ClickLogRepo
	var statusCheckRepo *postgres.StatusCheckRepo
	var guideRepo *postgres.GuideRepo
	var announcementRepo *postgres.AnnouncementRepo
	var feedbackRepo *postgres.FeedbackRepo
	var analyticsRepo *postgres.AnalyticsRepo
	var auditLogRepo *postgres.AuditLogRepo
	var rbacRepo *postgres.RBACRepo

	if err != nil {
		log.Warn().Err(err).Msg("Database connection failed in serverless init. Operating in offline DB mode.")
	} else {
		// Run migrations automatically
		if err := database.RunMigrations(context.Background(), db.Pool, "migrations"); err != nil {
			log.Warn().Err(err).Msg("Database migration check in serverless mode failed or skipped")
		}
		userRepo = postgres.NewUserRepo(db.Pool)
		categoryRepo = postgres.NewCategoryRepo(db.Pool)
		appRepo = postgres.NewAppRepo(db.Pool)
		favoriteRepo = postgres.NewFavoriteRepo(db.Pool)
		clickLogRepo = postgres.NewClickLogRepo(db.Pool)
		statusCheckRepo = postgres.NewStatusCheckRepo(db.Pool)
		guideRepo = postgres.NewGuideRepo(db.Pool)
		announcementRepo = postgres.NewAnnouncementRepo(db.Pool)
		feedbackRepo = postgres.NewFeedbackRepo(db.Pool)
		analyticsRepo = postgres.NewAnalyticsRepo(db.Pool)
		auditLogRepo = postgres.NewAuditLogRepo(db.Pool)
		rbacRepo = postgres.NewRBACRepo(db.Pool)
	}

	// 4. Initialize Background Workers & Services
	clickWorker := worker.NewClickWorker(clickLogRepo, 100, 1)
	sseHub := realtime.NewSSEHub()
	healthProbeWorker := worker.NewHealthProbeWorker(appRepo, statusCheckRepo, sseHub, 5*time.Minute)
	jwtService := jwt.NewJWTService(&cfg.JWT)
	ssoClient := sso.NewClient(&cfg.SSO)

	authUsecase := usecase.NewAuthUsecase(userRepo, ssoClient, jwtService)
	catalogUsecase := usecase.NewCatalogUsecase(categoryRepo, appRepo, userRepo, guideRepo)
	interactionUsecase := usecase.NewInteractionUsecase(favoriteRepo, clickLogRepo, appRepo, clickWorker)
	adminUsecase := usecase.NewAdminUsecase(appRepo, categoryRepo, guideRepo)
	monitoringUsecase := usecase.NewMonitoringUsecase(statusCheckRepo, appRepo, healthProbeWorker)
	announcementUsecase := usecase.NewAnnouncementUsecase(announcementRepo, appRepo)
	feedbackUsecase := usecase.NewFeedbackUsecase(feedbackRepo, appRepo)
	analyticsUsecase := usecase.NewAnalyticsUsecase(analyticsRepo)
	auditUsecase := usecase.NewAuditUsecase(auditLogRepo)
	rbacUsecase := usecase.NewRBACUsecase(rbacRepo)
	aiService := ai.NewAssistantService(appRepo, categoryRepo, guideRepo)
	aiUsecase := usecase.NewAIUsecase(aiService, userRepo)
	reportExporter := exporter.NewReportExporter()
	reportUsecase := usecase.NewReportUsecase(appRepo, analyticsRepo, reportExporter)
	imageOptimizer := media.NewImageOptimizer("/tmp/uploads/logos", "/uploads/logos")

	// 5. Setup Delivery & Handlers
	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, db)
	authHandler := deliveryHTTP.NewAuthHandler(authUsecase, cfg)
	catalogHandler := deliveryHTTP.NewCatalogHandler(catalogUsecase)
	interactionHandler := deliveryHTTP.NewInteractionHandler(interactionUsecase)
	adminHandler := deliveryHTTP.NewAdminHandler(adminUsecase, imageOptimizer)
	monitoringHandler := deliveryHTTP.NewMonitoringHandler(monitoringUsecase, sseHub)
	announcementHandler := deliveryHTTP.NewAnnouncementHandler(announcementUsecase)
	feedbackHandler := deliveryHTTP.NewFeedbackHandler(feedbackUsecase)
	analyticsHandler := deliveryHTTP.NewAnalyticsHandler(analyticsUsecase)
	auditHandler := deliveryHTTP.NewAuditHandler(auditUsecase)
	aiHandler := deliveryHTTP.NewAIHandler(aiUsecase)
	reportHandler := deliveryHTTP.NewReportHandler(reportUsecase)
	rbacHandler := deliveryHTTP.NewRBACHandler(rbacUsecase)

	appRouter = deliveryHTTP.SetupRouter(
		cfg,
		healthHandler,
		authHandler,
		catalogHandler,
		interactionHandler,
		adminHandler,
		monitoringHandler,
		announcementHandler,
		feedbackHandler,
		analyticsHandler,
		auditHandler,
		aiHandler,
		reportHandler,
		rbacHandler,
		jwtService,
	)
}

// Handler is the entrypoint for Vercel Serverless Functions
func Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Error().Interface("panic", rec).Msg("Serverless handler panic recovered")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success":false,"message":"Internal server error during serverless execution"}`))
		}
	}()

	initOnce.Do(initApp)

	if appRouter != nil {
		appRouter.ServeHTTP(w, r)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"message":"Server initialization failed"}`))
	}
}
