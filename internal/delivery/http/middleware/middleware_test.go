package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend-alusi-go/config"
	"backend-alusi-go/internal/delivery/http/middleware"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/pkg/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRequestIDAndSecurityHeadersMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.SecurityHeaders())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get(middleware.HeaderRequestID) == "" {
		t.Error("Expected X-Request-ID header to be set")
	}

	if w.Header().Get(middleware.HeaderResponseTime) == "" {
		t.Error("Expected X-Response-Time header to be set")
	}

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Expected X-Content-Type-Options 'nosniff', got '%s'", w.Header().Get("X-Content-Type-Options"))
	}

	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("Expected X-Frame-Options 'DENY', got '%s'", w.Header().Get("X-Frame-Options"))
	}
}

func TestCORSMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(middleware.CORS([]string{"http://localhost:3000"}))

	router.GET("/test-cors", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 1. Preflight OPTIONS request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/test-cors", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Expected status 204 for OPTIONS, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Expected Allow-Origin 'http://localhost:3000', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestRateLimiterMiddleware(t *testing.T) {
	limiter := middleware.NewIPRateLimiter(rate.Limit(2), 2, 1*time.Minute)

	router := gin.New()
	router.Use(middleware.RateLimit(limiter))

	router.GET("/limited", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Request 1 & 2 should succeed
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/limited", nil)
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("Request %d expected 200, got %d", i+1, w.Code)
		}
	}

	// Request 3 should be rate limited (429)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/limited", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("Request 3 expected 429 Too Many Requests, got %d", w.Code)
	}
}

func TestAuthAndRBACMiddleware(t *testing.T) {
	jwtCfg := &config.JWTConfig{
		Secret:          "test-secret-at-least-32-chars-long-for-testing",
		ExpirationHours: 24,
	}
	jwtService := jwt.NewJWTService(jwtCfg)

	// Users with different roles
	adminUser := &domain.User{
		ID:       1,
		SSOSub:   "sub_admin",
		Nama:     "Admin User",
		Email:    "admin@bps.go.id",
		Roles:    []domain.Role{{ID: 1, Nama: "admin"}},
		UserType: "internal",
	}

	pegawaiUser := &domain.User{
		ID:       2,
		SSOSub:   "sub_pegawai",
		Nama:     "Pegawai User",
		Email:    "pegawai@bps.go.id",
		Roles:    []domain.Role{{ID: 3, Nama: "pegawai"}},
		UserType: "internal",
	}

	adminToken, _, _ := jwtService.GenerateSessionToken(adminUser, "1200")
	pegawaiToken, _, _ := jwtService.GenerateSessionToken(pegawaiUser, "1200")

	router := gin.New()

	// Protected Admin endpoint
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthRequired(jwtService))
	adminGroup.Use(middleware.RequireRoles("admin"))
	adminGroup.GET("/dashboard", func(c *gin.Context) {
		c.String(http.StatusOK, "admin ok")
	})

	// 1. Unauthenticated -> 401
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got %d", w.Code)
	}

	// 2. Pegawai accessing Admin route -> 403 Forbidden
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+pegawaiToken)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden for non-admin, got %d", w.Code)
	}

	// 3. Admin accessing Admin route -> 200 OK
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for admin, got %d", w.Code)
	}
}
