package http

import (
	"fmt"
	"net/http"
	"time"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportUsecase *usecase.ReportUsecase
}

func NewReportHandler(reportUsecase *usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{
		reportUsecase: reportUsecase,
	}
}

// GetOpenAPICatalog godoc
// @Summary      Metadata Katalog Open API (Satu Data Indonesia / DCAT-AP)
// @Description  Mengambil katalog aplikasi dalam format standar metadata terbuka Satu Data Indonesia & DCAT-AP
// @Tags         Open Data / SDI
// @Produce      json
// @Success      200  {object}  usecase.OpenAPICatalogResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /openapi/apps [get]
func (h *ReportHandler) GetOpenAPICatalog(c *gin.Context) {
	data, err := h.reportUsecase.GetOpenAPICatalog(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal memuat metadata Open API: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// ExportAnalyticsCSV godoc
// @Summary      Ekspor laporan analitik (CSV / Excel)
// @Description  Mengunduh file laporan rekapitulasi statistik penggunaan dan SLA dalam format CSV (Admin)
// @Tags         Admin - Reports
// @Produce      text/csv
// @Security     BearerAuth
// @Success      200  {file}  binary
// @Failure      401  {object}  response.StandardResponse
// @Failure      403  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /admin/reports/analytics/export [get]
func (h *ReportHandler) ExportAnalyticsCSV(c *gin.Context) {
	bytesData, err := h.reportUsecase.GenerateAnalyticsCSV(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal mengekspor laporan analitik: "+err.Error())
		return
	}

	filename := fmt.Sprintf("laporan_analitik_portal_bps_%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", bytesData)
}

// ExportCatalogCSV godoc
// @Summary      Ekspor inventaris katalog aplikasi (CSV / Excel)
// @Description  Mengunduh file inventaris lengkap katalog aplikasi, URL, dan status dalam format CSV (Admin)
// @Tags         Admin - Reports
// @Produce      text/csv
// @Security     BearerAuth
// @Success      200  {file}  binary
// @Failure      401  {object}  response.StandardResponse
// @Failure      403  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /admin/reports/catalog/export [get]
func (h *ReportHandler) ExportCatalogCSV(c *gin.Context) {
	bytesData, err := h.reportUsecase.GenerateCatalogCSV(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal mengekspor katalog aplikasi: "+err.Error())
		return
	}

	filename := fmt.Sprintf("katalog_aplikasi_bps_%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", bytesData)
}
