package http

import (
	"net/http"
	"strconv"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsUsecase *usecase.AnalyticsUsecase
}

func NewAnalyticsHandler(analyticsUsecase *usecase.AnalyticsUsecase) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsUsecase: analyticsUsecase,
	}
}

// GetDashboardSummary returns high-level KPI rollups (DAU, MAU, Clicks, Totals)
func (h *AnalyticsHandler) GetDashboardSummary(c *gin.Context) {
	summary, err := h.analyticsUsecase.GetDashboardSummary(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal memuat ringkasan dasbor: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ringkasan metrik analitik berhasil dimuat", summary, nil)
}

// GetTopApps returns top application ranking by clicks
func (h *AnalyticsHandler) GetTopApps(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	topApps, err := h.analyticsUsecase.GetTopApps(c.Request.Context(), days, limit)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat ranking aplikasi: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ranking aplikasi terpopuler berhasil dimuat", topApps, nil)
}

// GetTrends returns daily click & user trends
func (h *AnalyticsHandler) GetTrends(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	trends, err := h.analyticsUsecase.GetDailyClicksTrend(c.Request.Context(), days)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat tren trafik: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Tren penggunaan aplikasi berhasil dimuat", trends, nil)
}

// GetDisruptions returns 30-day disruption and SLA uptime summary
func (h *AnalyticsHandler) GetDisruptions(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	disruptions, err := h.analyticsUsecase.GetDisruptionSummary(c.Request.Context(), days)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat ringkasan gangguan: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ringkasan gangguan dan SLA uptime berhasil dimuat", disruptions, nil)
}
