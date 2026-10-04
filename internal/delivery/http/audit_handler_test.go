package http_test

import (
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

func TestAuditEndpoints_Routing(t *testing.T) {
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
	analyticsUsecase := usecase.NewAnalyticsUsecase(nil)
	auditUsecase := usecase.NewAuditUsecase(nil)
	reportUsecase := usecase.NewReportUsecase(nil, nil, exporter.NewReportExporter())

	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, nil)
	authHandler := deliveryHTTP.NewAuthHandler(authUsecase, cfg)
	catalogHandler := deliveryHTTP.NewCatalogHandler(catalogUsecase)
	interactionHandler := deliveryHTTP.NewInteractionHandler(interactionUsecase)
	adminHandler := deliveryHTTP.NewAdminHandler(adminUsecase)
	monitoringHandler := deliveryHTTP.NewMonitoringHandler(monitoringUsecase, nil)
	announcementHandler := deliveryHTTP.NewAnnouncementHandler(announcementUsecase)
	feedbackHandler := deliveryHTTP.NewFeedbackHandler(feedbackUsecase)
	analyticsHandler := deliveryHTTP.NewAnalyticsHandler(analyticsUsecase)
	auditHandler := deliveryHTTP.NewAuditHandler(auditUsecase)
	aiUsecase := usecase.NewAIUsecase(ai.NewAssistantService(nil, nil, nil), nil)
	aiHandler := deliveryHTTP.NewAIHandler(aiUsecase)
	reportHandler := deliveryHTTP.NewReportHandler(reportUsecase)

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
		jwtService,
	)

	// Generate admin token
	adminUser := &domain.User{
		ID:       1,
		SSOSub:   "19950101",
		Nama:     "Admin BPS",
		Email:    "admin@bps.go.id",
		UserType: "internal",
		Roles: []domain.Role{
			{ID: 1, Nama: "admin"},
		},
	}
	adminToken, _, err := jwtService.GenerateSessionToken(adminUser, "1200")
	if err != nil {
		t.Fatalf("Failed to generate admin token: %v", err)
	}

	// 1. Test GET /api/v1/admin/audit-logs with Admin Token -> 200 OK
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/audit-logs", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for /admin/audit-logs, got %d", w.Code)
	}

	var resp response.StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success = true for /admin/audit-logs")
	}
}
