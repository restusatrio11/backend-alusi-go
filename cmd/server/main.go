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
	"backend-alusi-go/pkg/database"
	"backend-alusi-go/pkg/logger"

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
		Msg("Starting Portal BPS Sumut Backend API...")

	// 3. Initialize PostgreSQL Database Connection Pool
	ctx, cancelInit := context.WithTimeout(context.Background(), 15*time.Second)
	db, err := database.NewPostgresDB(ctx, &cfg.Database)
	cancelInit()

	if err != nil {
		log.Warn().Err(err).Msg("Database connection failed or not available yet. Server starting in offline DB mode.")
	} else {
		// Run migrations automatically
		if err := database.RunMigrations(context.Background(), db.Pool, "migrations"); err != nil {
			log.Error().Err(err).Msg("Database migration failed")
		} else {
			log.Info().Msg("All database migrations verified and applied")
		}
	}

	// 4. Setup Delivery & Handlers
	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, db)
	router := deliveryHTTP.SetupRouter(cfg, healthHandler)

	// 5. Configure HTTP Server
	serverAddr := fmt.Sprintf(":%s", cfg.App.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Start Server in Goroutine
	go func() {
		log.Info().Str("addr", serverAddr).Msg("HTTP Server is listening and serving requests")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("HTTP server failed to start")
		}
	}()

	// 7. Graceful Shutdown Listener
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Warn().Msg("Shutdown signal received, shutting down gracefully...")

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
