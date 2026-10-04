package http

import (
	"context"
	"net/http"
	"time"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/pkg/database"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	appName string
	appEnv  string
	db      *database.PostgresDB
}

func NewHealthHandler(appName, appEnv string, db *database.PostgresDB) *HealthHandler {
	return &HealthHandler{
		appName: appName,
		appEnv:  appEnv,
		db:      db,
	}
}

// Check handles /healthz endpoint
func (h *HealthHandler) Check(c *gin.Context) {
	dbStatus := "DISCONNECTED"
	if h.db != nil && h.db.Pool != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := h.db.Ping(ctx); err == nil {
			dbStatus = "CONNECTED"
		}
	}

	statusCode := http.StatusOK
	if h.db != nil && dbStatus != "CONNECTED" {
		statusCode = http.StatusServiceUnavailable
	}

	response.Success(c, statusCode, "Health status check", gin.H{
		"app":       h.appName,
		"env":       h.appEnv,
		"status":    "UP",
		"database":  dbStatus,
		"timestamp": time.Now().Format(time.RFC3339),
	}, nil)
}
