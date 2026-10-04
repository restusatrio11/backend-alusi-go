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

// GetDashboardSummary godoc
// @Summary      Ringkasan metrik eksekutif (KPI Rollup)
// @Description  Mengambil agregasi total pengguna, total aplikasi aktif, total klik, serta DAU dan MAU (Admin)
// @Tags         Admin - Analytics
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.StandardResponse
// @Failure      401  {object}  response.StandardResponse
// @Failure      403  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /admin/analytics/summary [get]
func (h *AnalyticsHandler) GetDashboardSummary(c *gin.Context) {
	summary, err := h.analyticsUsecase.GetDashboardSummary(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal memuat ringkasan dasbor: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ringkasan metrik analitik berhasil dimuat", summary, nil)
}

// GetTopApps godoc
// @Summary      Ranking aplikasi paling sering diakses
// @Description  Mengambil daftar peringkat aplikasi terpopuler berdasarkan jumlah klik (Admin)
// @Tags         Admin - Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        days   query     int  false  "Rentang hari ke belakang (default 30)"
// @Param        limit  query     int  false  "Jumlah item ranking (default 10)"
// @Success      200    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/analytics/top-apps [get]
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

// GetTrends godoc
// @Summary      Tren trafik dan aktivitas harian
// @Description  Mengambil data tren grafik klik harian dan jumlah pengguna aktif unik (Admin)
// @Tags         Admin - Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        days  query     int  false  "Rentang hari ke belakang (default 30)"
// @Success      200   {object}  response.StandardResponse
// @Failure      401   {object}  response.StandardResponse
// @Failure      403   {object}  response.StandardResponse
// @Failure      500   {object}  response.StandardResponse
// @Router       /admin/analytics/trends [get]
func (h *AnalyticsHandler) GetTrends(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	trends, err := h.analyticsUsecase.GetDailyClicksTrend(c.Request.Context(), days)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat tren trafik: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Tren penggunaan aplikasi berhasil dimuat", trends, nil)
}

// GetDisruptions godoc
// @Summary      Rekapitulasi gangguan layanan dan SLA
// @Description  Mengambil riwayat insiden kendala operasional dan persentase uptime SLA per layanan (Admin)
// @Tags         Admin - Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        days  query     int  false  "Rentang hari ke belakang (default 30)"
// @Success      200   {object}  response.StandardResponse
// @Failure      401   {object}  response.StandardResponse
// @Failure      403   {object}  response.StandardResponse
// @Failure      500   {object}  response.StandardResponse
// @Router       /admin/analytics/disruptions [get]
func (h *AnalyticsHandler) GetDisruptions(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	disruptions, err := h.analyticsUsecase.GetDisruptionSummary(c.Request.Context(), days)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat ringkasan gangguan: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ringkasan gangguan dan SLA uptime berhasil dimuat", disruptions, nil)
}
