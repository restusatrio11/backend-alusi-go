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

// GetOpenAPICatalog returns Dublin Core / SDI standard Open API metadata
func (h *ReportHandler) GetOpenAPICatalog(c *gin.Context) {
	data, err := h.reportUsecase.GetOpenAPICatalog(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal memuat metadata Open API: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// ExportAnalyticsCSV exports analytics summary to downloadable CSV file
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

// ExportCatalogCSV exports application catalog to downloadable CSV file
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
