package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-alusi-go/config"
	deliveryHTTP "backend-alusi-go/internal/delivery/http"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/ai"
	"backend-alusi-go/pkg/exporter"
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/sso"
)

func TestReportEndpoints_Routing(t *testing.T) {
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
	aiUsecase := usecase.NewAIUsecase(ai.NewAssistantService(nil, nil, nil), nil)

	reportExporter := exporter.NewReportExporter()
	reportUsecase := usecase.NewReportUsecase(nil, nil, reportExporter)

	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, nil)
	authHandler := deliveryHTTP.NewAuthHandler(authUsecase, cfg)
	catalogHandler := deliveryHTTP.NewCatalogHandler(catalogUsecase)
	interactionHandler := deliveryHTTP.NewInteractionHandler(interactionUsecase)
	adminHandler := deliveryHTTP.NewAdminHandler(adminUsecase, nil)
	monitoringHandler := deliveryHTTP.NewMonitoringHandler(monitoringUsecase, nil)
	announcementHandler := deliveryHTTP.NewAnnouncementHandler(announcementUsecase)
	feedbackHandler := deliveryHTTP.NewFeedbackHandler(feedbackUsecase)
	analyticsHandler := deliveryHTTP.NewAnalyticsHandler(analyticsUsecase)
	auditHandler := deliveryHTTP.NewAuditHandler(auditUsecase)
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

	// 1. Test GET /api/v1/openapi/apps -> 200 OK
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/openapi/apps", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for /openapi/apps, got %d", w.Code)
	}

	// 2. Test Admin Export with Admin Token
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
	adminToken, _, _ := jwtService.GenerateSessionToken(adminUser, "1200")

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/reports/catalog/export", nil)
	req2.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for /admin/reports/catalog/export, got %d", w2.Code)
	}
}
