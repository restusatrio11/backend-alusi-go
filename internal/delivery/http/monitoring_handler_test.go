package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-alusi-go/config"
	deliveryHTTP "backend-alusi-go/internal/delivery/http"
	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/sso"
)

func TestMonitoringEndpoints_Routing(t *testing.T) {
	cfg := &config.Config{
		App:  config.AppConfig{Name: "test-app", Env: "test", Debug: true},
		JWT:  config.JWTConfig{Secret: "test-secret-at-least-32-chars-long", ExpirationHours: 24},
		SSO:  config.SSOConfig{ClientID: "test-client", IssuerURL: "https://otp-dev.bps.web.id"},
		CORS: config.CORSConfig{AllowedOrigins: []string{"*"}},
	}

	jwtService := jwt.NewJWTService(&cfg.JWT)
	ssoClient := sso.NewClient(&cfg.SSO)
	authUsecase := usecase.NewAuthUsecase(nil, ssoClient, jwtService)
	catalogUsecase := usecase.NewCatalogUsecase(nil, nil, nil)
	interactionUsecase := usecase.NewInteractionUsecase(nil, nil, nil, nil)
	adminUsecase := usecase.NewAdminUsecase(nil, nil)
	monitoringUsecase := usecase.NewMonitoringUsecase(nil, nil, nil)

	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, nil)
	authHandler := deliveryHTTP.NewAuthHandler(authUsecase, cfg)
	catalogHandler := deliveryHTTP.NewCatalogHandler(catalogUsecase)
	interactionHandler := deliveryHTTP.NewInteractionHandler(interactionUsecase)
	adminHandler := deliveryHTTP.NewAdminHandler(adminUsecase)
	monitoringHandler := deliveryHTTP.NewMonitoringHandler(monitoringUsecase)

	router := deliveryHTTP.SetupRouter(
		cfg,
		healthHandler,
		authHandler,
		catalogHandler,
		interactionHandler,
		adminHandler,
		monitoringHandler,
		jwtService,
	)

	// 1. Test GET /api/v1/services/status -> 200 OK
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/services/status", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for /services/status, got %d", w.Code)
	}

	var resp response.StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success = true for /services/status")
	}

	// 2. Test GET /api/v1/apps/simbatik/status-history -> 200 OK (empty list when offline)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/apps/simbatik/status-history", nil)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for /apps/simbatik/status-history, got %d", w2.Code)
	}
}
