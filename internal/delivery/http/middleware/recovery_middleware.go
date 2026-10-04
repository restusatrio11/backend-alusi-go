package middleware

import (
	"fmt"
	"net/http"

	"backend-alusi-go/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Recovery middleware recovers from any panics and writes a 500 error response
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				reqID, _ := c.Get("request_id")
				log.Error().
					Str("request_id", fmt.Sprintf("%v", reqID)).
					Interface("panic", err).
					Msg("Unhandled panic recovered")

				response.Error(
					c,
					http.StatusInternalServerError,
					"INTERNAL_SERVER_ERROR",
					"Terjadi kesalahan internal server.",
					nil,
				)
				c.Abort()
			}
		}()
		c.Next()
	}
}
