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

	// 3. Setup Delivery & Handlers
	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env)
	router := deliveryHTTP.SetupRouter(cfg, healthHandler)

	// 4. Configure HTTP Server
	serverAddr := fmt.Sprintf(":%s", cfg.App.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Start Server in Goroutine
	go func() {
		log.Info().Str("addr", serverAddr).Msg("HTTP Server is listening and serving requests")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("HTTP server failed to start")
		}
	}()

	// 6. Graceful Shutdown Listener
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Warn().Msg("Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		log.Info().Msg("Server exiting cleanly")
	}
}
