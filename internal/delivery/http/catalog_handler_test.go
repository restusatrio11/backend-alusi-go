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
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/sso"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCatalogEndpoints_Routing(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{Name: "test-app", Env: "test", Debug: true},
		JWT: config.JWTConfig{Secret: "test-secret-at-least-32-chars-long", ExpirationHours: 24},
		SSO: config.SSOConfig{ClientID: "test-client", IssuerURL: "https://otp-dev.bps.web.id"},
		CORS: config.CORSConfig{AllowedOrigins: []string{"*"}},
	}

	jwtService := jwt.NewJWTService(&cfg.JWT)
	ssoClient := sso.NewClient(&cfg.SSO)
	authUsecase := usecase.NewAuthUsecase(nil, ssoClient, jwtService)
	catalogUsecase := usecase.NewCatalogUsecase(nil, nil, nil)

	healthHandler := deliveryHTTP.NewHealthHandler(cfg.App.Name, cfg.App.Env, nil)
	authHandler := deliveryHTTP.NewAuthHandler(authUsecase, cfg)
	catalogHandler := deliveryHTTP.NewCatalogHandler(catalogUsecase)

	router := deliveryHTTP.SetupRouter(cfg, healthHandler, authHandler, catalogHandler, jwtService)

	// 1. Test GET /api/v1/apps/search with empty query -> 400 Bad Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/apps/search", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for empty search query, got %d", w.Code)
	}

	var resp response.StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if resp.Success {
		t.Errorf("Expected success = false for bad request")
	}
	if resp.Error == nil || resp.Error.Code != "BAD_REQUEST" {
		t.Errorf("Expected error code BAD_REQUEST, got %v", resp.Error)
	}
}

func TestDomainEntities(t *testing.T) {
	cat := domain.Category{
		ID:   1,
		Nama: "Survei & Sensus",
		Slug: "survei-sensus",
	}

	app := domain.App{
		ID:             1,
		CategoryID:     cat.ID,
		Category:       &cat,
		Nama:           "SIMBATIK",
		Slug:           "simbatik",
		URL:            "https://simbatik.bps.go.id",
		TargetPengguna: "Petugas Lapangan",
		StatusLayanan:  "online",
		IsPublic:       true,
		Aktif:          true,
	}

	if app.Category.Nama != "Survei & Sensus" {
		t.Errorf("Expected category name 'Survei & Sensus', got '%s'", app.Category.Nama)
	}
	if !app.IsPublic || !app.Aktif {
		t.Errorf("Expected app to be active and public")
	}
}
