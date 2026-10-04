package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger middleware logs incoming HTTP requests with Zerolog
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		reqID, _ := c.Get("request_id")

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		event := log.Info()
		if status >= 400 && status < 500 {
			event = log.Warn()
		} else if status >= 500 {
			event = log.Error()
		}

		event.
			Str("request_id", reqID.(string)).
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Str("client_ip", clientIP).
			Msg("HTTP Request")
	}
}
