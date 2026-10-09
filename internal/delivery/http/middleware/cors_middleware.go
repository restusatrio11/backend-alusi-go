package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigin := ""

		if origin == "" {
			allowOrigin = "*"
		} else {
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowOrigin = origin
					break
				}
			}
			// Fallback: If allowedOrigins has wildcard "*" or matches local LAN subnet
			if allowOrigin == "" {
				for _, o := range allowedOrigins {
					if o == "*" {
						allowOrigin = origin
						break
					}
				}
			}
			if allowOrigin == "" {
				if strings.HasPrefix(origin, "http://10.") ||
					strings.HasPrefix(origin, "http://192.168.") ||
					strings.HasPrefix(origin, "http://172.") ||
					strings.HasPrefix(origin, "http://localhost") ||
					strings.HasPrefix(origin, "http://127.0.0.1") {
					allowOrigin = origin
				}
			}
		}

		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
