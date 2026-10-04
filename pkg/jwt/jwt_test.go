package jwt_test

import (
	"testing"

	"backend-alusi-go/config"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/pkg/jwt"
)

func TestJWTService(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:          "test-super-secret-key-at-least-32-chars-long",
		ExpirationHours: 24,
	}

	jwtService := jwt.NewJWTService(cfg)
	nip := "199501012020011001"
	user := &domain.User{
		ID:       10,
		SSOSub:   "user_12345",
		UserType: "internal",
		Nama:     "Restu Satrio Pinanggih",
		Email:    "restu@bps.go.id",
		NIP:      &nip,
		Roles: []domain.Role{
			{ID: 1, Nama: "admin"},
			{ID: 3, Nama: "pegawai"},
		},
	}

	token, expiresAt, err := jwtService.GenerateSessionToken(user, "1200")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Expected non-empty token string")
	}

	if expiresAt.IsZero() {
		t.Fatal("Expected valid expiration time")
	}

	claims, err := jwtService.ValidateSessionToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != 10 {
		t.Errorf("Expected UserID 10, got %d", claims.UserID)
	}
	if claims.SSOSub != "user_12345" {
		t.Errorf("Expected SSOSub 'user_12345', got '%s'", claims.SSOSub)
	}
	if claims.UserType != "internal" {
		t.Errorf("Expected UserType 'internal', got '%s'", claims.UserType)
	}
	if claims.SatkerKode != "1200" {
		t.Errorf("Expected SatkerKode '1200', got '%s'", claims.SatkerKode)
	}
	if len(claims.Roles) != 2 || claims.Roles[0] != "admin" {
		t.Errorf("Unexpected roles: %v", claims.Roles)
	}
}
