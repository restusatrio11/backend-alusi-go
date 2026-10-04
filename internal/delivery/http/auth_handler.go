package http

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"backend-alusi-go/config"
	"backend-alusi-go/internal/delivery/http/middleware"
	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const CookieOAuthState = "alusi_oauth_state"

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	cfg         *config.Config
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		cfg:         cfg,
	}
}

// Login handles initiation of SSO authentication
func (h *AuthHandler) Login(c *gin.Context) {
	// Generate random 16-byte state string
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	state := hex.EncodeToString(b)

	// Store state in short-lived cookie for validation in callback
	c.SetCookie(
		CookieOAuthState,
		state,
		300, // 5 minutes
		"/",
		h.cfg.JWT.CookieDomain,
		h.cfg.JWT.CookieSecure,
		true, // HttpOnly
	)

	authURL := h.authUsecase.GetLoginURL(state)

	// If query has redirect=false, return JSON url instead of 302 redirect
	if c.Query("redirect") == "false" {
		response.Success(c, http.StatusOK, "Login URL generated", gin.H{
			"login_url": authURL,
			"state":     state,
		}, nil)
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

// Callback handles authorization code exchange from SSO Sumut
func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.BadRequest(c, "Authorization code tidak ditemukan pada parameter query.", nil)
		return
	}

	state := c.Query("state")
	cookieState, _ := c.Cookie(CookieOAuthState)

	// Clean up state cookie
	c.SetCookie(CookieOAuthState, "", -1, "/", h.cfg.JWT.CookieDomain, h.cfg.JWT.CookieSecure, true)

	if state == "" || cookieState == "" || state != cookieState {
		log.Warn().Str("state", state).Str("cookie_state", cookieState).Msg("OAuth state mismatch or expired")
		// Optional: allow proceed in development mode if debug is enabled
		if !h.cfg.App.Debug {
			response.BadRequest(c, "OAuth state tidak valid atau sesi login telah kedaluwarsa.", nil)
			return
		}
	}

	user, sessionToken, err := h.authUsecase.HandleSSOCallback(c.Request.Context(), code)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process SSO callback")
		response.InternalServerError(c, "Gagal memproses autentikasi SSO BPS: "+err.Error())
		return
	}

	// Set session cookie
	maxAge := h.cfg.JWT.ExpirationHours * 3600
	c.SetCookie(
		middleware.CookieSessionName,
		sessionToken,
		maxAge,
		"/",
		h.cfg.JWT.CookieDomain,
		h.cfg.JWT.CookieSecure,
		true, // HttpOnly
	)

	// If request came from API/XHR or accept JSON
	acceptHeader := c.GetHeader("Accept")
	if strings.Contains(acceptHeader, "application/json") || c.Query("format") == "json" {
		response.Success(c, http.StatusOK, "Login berhasil", gin.H{
			"token": sessionToken,
			"user":  user,
		}, nil)
		return
	}

	// Otherwise, redirect to frontend homepage / portal dashboard
	frontendOrigin := "http://localhost:3000"
	if len(h.cfg.CORS.AllowedOrigins) > 0 && h.cfg.CORS.AllowedOrigins[0] != "*" {
		frontendOrigin = h.cfg.CORS.AllowedOrigins[0]
	}
	c.Redirect(http.StatusFound, frontendOrigin+"?login=success")
}

// Me returns currently authenticated user details
func (h *AuthHandler) Me(c *gin.Context) {
	claimsVal, exists := c.Get(middleware.CtxUserClaims)
	if !exists {
		response.Unauthorized(c, "Pengguna belum terautentikasi.")
		return
	}

	claims := claimsVal.(*jwt.SessionClaims)
	user, err := h.authUsecase.GetProfile(c.Request.Context(), claims.UserID)
	if err != nil || user == nil {
		response.NotFound(c, "Data pengguna tidak ditemukan.")
		return
	}

	response.Success(c, http.StatusOK, "Profil pengguna berhasil dimuat", user, nil)
}

// Logout handles Single Log Out (SLO)
func (h *AuthHandler) Logout(c *gin.Context) {
	// 1. Clear local application session cookie
	c.SetCookie(
		middleware.CookieSessionName,
		"",
		-1,
		"/",
		h.cfg.JWT.CookieDomain,
		h.cfg.JWT.CookieSecure,
		true,
	)

	postLogoutRedirect := c.Query("post_logout_redirect_uri")
	logoutURL := h.authUsecase.GetLogoutURL(postLogoutRedirect)

	if c.Request.Method == http.MethodPost || c.Query("redirect") == "false" {
		response.Success(c, http.StatusOK, "Logout berhasil", gin.H{
			"logout_url": logoutURL,
		}, nil)
		return
	}

	c.Redirect(http.StatusFound, logoutURL)
}
