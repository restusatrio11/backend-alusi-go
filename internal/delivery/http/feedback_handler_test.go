package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-alusi-go/config"
	deliveryHTTP "backend-alusi-go/internal/delivery/http"
	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/ai"
	"backend-alusi-go/pkg/exporter"
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/sso"
)

func TestFeedbackEndpoints_Routing(t *testing.T) {
	cfg := &config.Config{
		App:  config.AppConfig{Name: "test-app", Env: "test", Debug: true},
		JWT:  config.JWTConfig{Secret: "test-secret-at-least-32-chars-long", ExpirationHours: 24},
		SSO:  config.SSOConfig{ClientID: "test-client", IssuerURL: "https://otp-dev.bps.web.id"},
		CORS: config.CORSConfig{AllowedOrigins: []string{"*"}},
	}

	jwtService := jwt.NewJWTService(&cfg.JWT)
	ssoClient := sso.NewClient(&cfg.SSO)
	authUsecase := usecase.NewAuthUsecase(nil, ssoClient, jwtService)
	catalogUsecase := usecase.NewCatalogUsecase(nil, nil, nil, nil)
	interactionUsecase := usecase.NewInteractionUsecase(nil, nil, nil, nil)
	adminUsecase := usecase.NewAdminUsecase(nil, nil, nil)
	monitoringUsecase := usecase.NewMonitoringUsecase(nil, nil, nil)
	announcementUsecase := usecase.NewAnnouncementUsecase(nil, nil)
	feedbackUsecase := usecase.NewFeedbackUsecase(nil, nil)
	reportUsecase := usecase.NewReportUsecase(nil, nil, exporter.NewReportExporter())

	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, nil)
	authHandler := deliveryHTTP.NewAuthHandler(authUsecase, cfg)
	catalogHandler := deliveryHTTP.NewCatalogHandler(catalogUsecase)
	interactionHandler := deliveryHTTP.NewInteractionHandler(interactionUsecase)
	adminHandler := deliveryHTTP.NewAdminHandler(adminUsecase, nil)
	monitoringHandler := deliveryHTTP.NewMonitoringHandler(monitoringUsecase, nil)
	announcementHandler := deliveryHTTP.NewAnnouncementHandler(announcementUsecase)
	feedbackHandler := deliveryHTTP.NewFeedbackHandler(feedbackUsecase)
	analyticsUsecase := usecase.NewAnalyticsUsecase(nil)
	analyticsHandler := deliveryHTTP.NewAnalyticsHandler(analyticsUsecase)
	auditUsecase := usecase.NewAuditUsecase(nil)
	auditHandler := deliveryHTTP.NewAuditHandler(auditUsecase)
	aiUsecase := usecase.NewAIUsecase(ai.NewAssistantService(nil, nil, nil), nil)
	aiHandler := deliveryHTTP.NewAIHandler(aiUsecase)
	reportHandler := deliveryHTTP.NewReportHandler(reportUsecase)
	rbacHandler := deliveryHTTP.NewRBACHandler(nil)

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
		auditHandler,
		aiHandler,
		reportHandler,
		rbacHandler,
		jwtService,
	)

	// 1. Test POST /api/v1/feedbacks -> 201 Created
	payload := map[string]interface{}{
		"kategori": "kendala",
		"pesan":    "Tidak bisa mengakses halaman login SIMBATIK",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/feedbacks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for POST /feedbacks, got %d", w.Code)
	}

	var resp response.StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success = true for POST /feedbacks")
	}
}

func TestFeedbackDomain(t *testing.T) {
	fb := domain.Feedback{
		ID:       1,
		Kategori: "kendala",
		Pesan:    "Error loading map",
		Status:   "pending",
	}

	if fb.Kategori != "kendala" || fb.Status != "pending" {
		t.Errorf("Unexpected feedback domain data")
	}
}
