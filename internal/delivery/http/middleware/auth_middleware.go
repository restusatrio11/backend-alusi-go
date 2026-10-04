package middleware

import (
	"strings"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	CookieSessionName = "alusi_session"
	CtxUserClaims     = "user_claims"
	CtxUserID         = "user_id"
	CtxUserRoles      = "user_roles"
)

// AuthRequired middleware enforces JWT authentication
func AuthRequired(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		// 1. Check Authorization Bearer Header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Fallback to Cookie
		if tokenString == "" {
			if cookie, err := c.Cookie(CookieSessionName); err == nil {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			response.Unauthorized(c, "Autentikasi diperlukan untuk mengakses resource ini.")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateSessionToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Sesi login Anda tidak valid atau telah kedaluwarsa.")
			c.Abort()
			return
		}

		c.Set(CtxUserClaims, claims)
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUserRoles, claims.Roles)

		c.Next()
	}
}

// OptionalAuth middleware extracts JWT authentication if present without blocking unauthenticated requests
func OptionalAuth(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if tokenString == "" {
			if cookie, err := c.Cookie(CookieSessionName); err == nil {
				tokenString = cookie
			}
		}

		if tokenString != "" {
			if claims, err := jwtService.ValidateSessionToken(tokenString); err == nil {
				c.Set(CtxUserClaims, claims)
				c.Set(CtxUserID, claims.UserID)
				c.Set(CtxUserRoles, claims.Roles)
			}
		}

		c.Next()
	}
}

// RequireRoles middleware restricts endpoint to specific roles
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesVal, exists := c.Get(CtxUserRoles)
		if !exists {
			response.Forbidden(c, "Anda tidak memiliki izin untuk mengakses resource ini.")
			c.Abort()
			return
		}

		userRoles, ok := rolesVal.([]string)
		if !ok {
			response.Forbidden(c, "Hak akses tidak valid.")
			c.Abort()
			return
		}

		hasRole := false
		for _, userRole := range userRoles {
			for _, allowed := range allowedRoles {
				if userRole == allowed {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "Hak akses Anda tidak mencukupi untuk aksi ini.")
			c.Abort()
			return
		}

		c.Next()
	}
}
