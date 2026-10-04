package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders middleware injects standard OWASP defensive HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME-type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking by restricting framing
		c.Header("X-Frame-Options", "DENY")

		// XSS Filter protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Control referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:;")

		// Restrict dangerous browser features
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		c.Next()
	}
}
