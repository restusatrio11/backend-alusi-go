package jwt

import (
	"errors"
	"fmt"
	"time"

	"backend-alusi-go/config"
	"backend-alusi-go/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type SessionClaims struct {
	UserID     int      `json:"user_id"`
	SSOSub     string   `json:"sso_sub"`
	UserType   string   `json:"user_type"`
	Nama       string   `json:"nama"`
	Email      string   `json:"email"`
	NIP        string   `json:"nip,omitempty"`
	SatkerKode string   `json:"satker_kode,omitempty"`
	Roles      []string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTService struct {
	cfg *config.JWTConfig
}

func NewJWTService(cfg *config.JWTConfig) *JWTService {
	return &JWTService{cfg: cfg}
}

// GenerateSessionToken creates a signed JWT session token for an authenticated user
func (s *JWTService) GenerateSessionToken(user *domain.User, satkerKode string) (string, time.Time, error) {
	expDuration := time.Duration(s.cfg.ExpirationHours) * time.Hour
	expiresAt := time.Now().Add(expDuration)

	var roleNames []string
	for _, r := range user.Roles {
		roleNames = append(roleNames, r.Nama)
	}

	nip := ""
	if user.NIP != nil {
		nip = *user.NIP
	}

	claims := SessionClaims{
		UserID:     user.ID,
		SSOSub:     user.SSOSub,
		UserType:   user.UserType,
		Nama:       user.Nama,
		Email:      user.Email,
		NIP:        nip,
		SatkerKode: satkerKode,
		Roles:      roleNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "alusi-portal-sumut",
			Subject:   user.SSOSub,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, expiresAt, nil
}

// ValidateSessionToken parses and validates a signed JWT token
func (s *JWTService) ValidateSessionToken(tokenString string) (*SessionClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*SessionClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
