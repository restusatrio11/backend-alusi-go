package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	HeaderRequestID    = "X-Request-ID"
	HeaderResponseTime = "X-Response-Time"
)

// RequestID middleware generates or propagates X-Request-ID and measures execution time
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := c.GetHeader(HeaderRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set("request_id", reqID)
		c.Header(HeaderRequestID, reqID)

		c.Next()

		duration := time.Since(start)
		c.Header(HeaderResponseTime, fmt.Sprintf("%dms", duration.Milliseconds()))
	}
}
