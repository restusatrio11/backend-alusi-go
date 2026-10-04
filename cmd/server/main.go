package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-alusi-go/config"
	deliveryHTTP "backend-alusi-go/internal/delivery/http"
	"backend-alusi-go/internal/repository/postgres"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/database"
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/logger"
	"backend-alusi-go/pkg/sso"
	"backend-alusi-go/pkg/worker"

	"github.com/rs/zerolog/log"
)

func main() {
	// 1. Load Configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Structured Logger
	logger.InitLogger(cfg.App.Debug, cfg.App.Env)
	log.Info().
		Str("app", cfg.App.Name).
		Str("env", cfg.App.Env).
		Str("port", cfg.App.Port).
		Str("sso_issuer", cfg.SSO.IssuerURL).
		Msg("Starting Portal BPS Sumut Backend API...")

	// 3. Initialize PostgreSQL Database Connection Pool
	ctx, cancelInit := context.WithTimeout(context.Background(), 15*time.Second)
	db, err := database.NewPostgresDB(ctx, &cfg.Database)
	cancelInit()

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

	if err != nil {
		log.Warn().Err(err).Msg("Database connection failed or not available yet. Server starting in offline DB mode.")
	} else {
		// Run migrations automatically
		if err := database.RunMigrations(context.Background(), db.Pool, "migrations"); err != nil {
			log.Error().Err(err).Msg("Database migration failed")
		} else {
			log.Info().Msg("All database migrations verified and applied")
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
	}

	// 4. Initialize Background Workers & Services
	clickWorker := worker.NewClickWorker(clickLogRepo, 1000, 3)
	healthProbeWorker := worker.NewHealthProbeWorker(appRepo, statusCheckRepo, 5*time.Minute)
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

	// 5. Setup Delivery & Handlers
	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, db)
	authHandler := deliveryHTTP.NewAuthHandler(authUsecase, cfg)
	catalogHandler := deliveryHTTP.NewCatalogHandler(catalogUsecase)
	interactionHandler := deliveryHTTP.NewInteractionHandler(interactionUsecase)
	adminHandler := deliveryHTTP.NewAdminHandler(adminUsecase)
	monitoringHandler := deliveryHTTP.NewMonitoringHandler(monitoringUsecase)
	announcementHandler := deliveryHTTP.NewAnnouncementHandler(announcementUsecase)
	feedbackHandler := deliveryHTTP.NewFeedbackHandler(feedbackUsecase)
	analyticsHandler := deliveryHTTP.NewAnalyticsHandler(analyticsUsecase)

	router := deliveryHTTP.SetupRouter(
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
		jwtService,
	)

	// Start health probe worker
	healthProbeWorker.Start()

	// 6. Configure HTTP Server
	serverAddr := fmt.Sprintf(":%s", cfg.App.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 7. Start Server in Goroutine
	go func() {
		log.Info().Str("addr", serverAddr).Msg("HTTP Server is listening and serving requests")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("HTTP server failed to start")
		}
	}()

	// 8. Graceful Shutdown Listener
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Warn().Msg("Shutdown signal received, shutting down gracefully...")

	// Stop background workers
	clickWorker.Stop()
	healthProbeWorker.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		log.Info().Msg("HTTP Server shut down cleanly")
	}

	// Close database connection pool
	if db != nil {
		db.Close()
	}
}
