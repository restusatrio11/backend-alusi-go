package http_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"backend-alusi-go/config"
	deliveryHTTP "backend-alusi-go/internal/delivery/http"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/ai"
	"backend-alusi-go/pkg/exporter"
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/media"
	"backend-alusi-go/pkg/realtime"
	"backend-alusi-go/pkg/sso"
)

func createTestPNGBytes(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 150, B: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func TestLogoUploadEndpoints_AuthAndStatic(t *testing.T) {
	tempUploadDir := filepath.Join(os.TempDir(), "alusi_http_test_uploads")
	_ = os.MkdirAll(tempUploadDir, 0755)
	defer os.RemoveAll(tempUploadDir)

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
	reportUsecase := usecase.NewReportUsecase(nil, nil, exporter.NewReportExporter())

	imageOptimizer := media.NewImageOptimizer(tempUploadDir, "/uploads/logos")
	sseHub := realtime.NewSSEHub()
	sseHub.Start()
	defer sseHub.Stop()

	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, nil)
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

	// 1. Test POST /api/v1/admin/apps/1/logo without Auth -> 401 Unauthorized
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/apps/1/logo", nil)
	req1.RequestURI = "/api/v1/admin/apps/1/logo"
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got %d", w1.Code)
	}

	// 2. Test POST /api/v1/admin/apps/1/logo with regular user role -> 403 Forbidden
	normalUser := &domain.User{
		ID:       10,
		Email:    "user@bps.go.id",
		Nama:     "User Test",
		UserType: "internal",
		Roles:    []domain.Role{{ID: 2, Nama: "user"}},
	}
	userToken, _, _ := jwtService.GenerateSessionToken(normalUser, "1200")
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/apps/1/logo", nil)
	req2.RequestURI = "/api/v1/admin/apps/1/logo"
	req2.Header.Set("Authorization", "Bearer "+userToken)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden for non-admin user, got %d", w2.Code)
	}

	// 3. Test POST /api/v1/admin/apps/1/logo with Admin token but missing file / uninitialized db -> 404 (App not found)
	adminUser := &domain.User{
		ID:       1,
		Email:    "admin@bps.go.id",
		Nama:     "Super Admin",
		UserType: "internal",
		Roles:    []domain.Role{{ID: 1, Nama: "admin"}},
	}
	adminToken, _, _ := jwtService.GenerateSessionToken(adminUser, "1200")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("logo", "logo.png")
	_, _ = part.Write(createTestPNGBytes(200, 200))
	_ = writer.Close()

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/apps/1/logo", body)
	req3.RequestURI = "/api/v1/admin/apps/1/logo"
	req3.Header.Set("Authorization", "Bearer "+adminToken)
	req3.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w3, req3)

	// Since appRepo is nil in unit mock, it returns 404 (Aplikasi tidak ditemukan)
	if w3.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 Not Found when app is missing in repo, got %d", w3.Code)
	}
}
