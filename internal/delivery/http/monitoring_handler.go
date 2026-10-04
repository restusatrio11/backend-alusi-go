package http

import (
	"net/http"
	"strconv"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type MonitoringHandler struct {
	monitoringUsecase *usecase.MonitoringUsecase
}

func NewMonitoringHandler(monitoringUsecase *usecase.MonitoringUsecase) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringUsecase: monitoringUsecase,
	}
}

// GetAppStatusHistory returns health check log history for an app
func (h *MonitoringHandler) GetAppStatusHistory(c *gin.Context) {
	slugOrID := c.Param("slug")
	if slugOrID == "" {
		slugOrID = c.Param("id")
	}
	if slugOrID == "" {
		response.BadRequest(c, "Identifier aplikasi tidak valid.", nil)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	history, err := h.monitoringUsecase.GetAppStatusHistoryByIdentifier(c.Request.Context(), slugOrID, limit)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Riwayat status aplikasi berhasil dimuat", history, nil)
}

// GetServiceUptimeSummary returns 30-day uptime summary for all services
func (h *MonitoringHandler) GetServiceUptimeSummary(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	summary, err := h.monitoringUsecase.GetServiceUptimeSummary(c.Request.Context(), days)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat ringkasan uptime layanan: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ringkasan uptime layanan berhasil dimuat", summary, nil)
}

// ManualProbeApp triggers an on-demand probe by admin
func (h *MonitoringHandler) ManualProbeApp(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID aplikasi tidak valid.", nil)
		return
	}

	check, err := h.monitoringUsecase.ManualProbeApp(c.Request.Context(), appID)
	if err != nil {
		response.InternalServerError(c, "Gagal melakukan health probe: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Health probe selesai dijalankan", check, nil)
}
