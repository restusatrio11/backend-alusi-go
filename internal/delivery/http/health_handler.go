package http

import (
	"net/http"
	"time"

	"backend-alusi-go/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	appName string
	appEnv  string
}

func NewHealthHandler(appName, appEnv string) *HealthHandler {
	return &HealthHandler{
		appName: appName,
		appEnv:  appEnv,
	}
}

// Check handles /healthz endpoint
func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, http.StatusOK, "Service is healthy", gin.H{
		"app":       h.appName,
		"env":       h.appEnv,
		"status":    "UP",
		"timestamp": time.Now().Format(time.RFC3339),
	}, nil)
}
