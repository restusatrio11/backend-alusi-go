package http

import (
	"net/http"
	"strconv"
	"time"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditUsecase *usecase.AuditUsecase
}

func NewAuditHandler(auditUsecase *usecase.AuditUsecase) *AuditHandler {
	return &AuditHandler{
		auditUsecase: auditUsecase,
	}
}

// ListAuditLogs godoc
// @Summary      Jejak aktivitas admin (Audit Trail)
// @Description  Mengambil log aktivitas operasional, perubahan data, dan akses administrator (Admin)
// @Tags         Admin - Audit Trail
// @Produce      json
// @Security     BearerAuth
// @Param        user_id     query     int     false  "Filter ID pengguna/admin"
// @Param        action      query     string  false  "Filter aksi (CREATE, UPDATE, DELETE, PROBE)"
// @Param        entity      query     string  false  "Filter entitas target (app, category, guide, announcement)"
// @Param        entity_id   query     string  false  "Filter ID entitas spesifik"
// @Param        start_date  query     string  false  "Filter tanggal awal (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "Filter tanggal akhir (YYYY-MM-DD)"
// @Param        page        query     int     false  "Halaman (default 1)"
// @Param        per_page    query     int     false  "Jumlah item per halaman (default 20)"
// @Success      200         {object}  response.StandardResponse
// @Failure      401         {object}  response.StandardResponse
// @Failure      403         {object}  response.StandardResponse
// @Failure      500         {object}  response.StandardResponse
// @Router       /admin/audit-logs [get]
func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	var userID *int
	if rawUID := c.Query("user_id"); rawUID != "" {
		if id, err := strconv.Atoi(rawUID); err == nil {
			userID = &id
		}
	}

	var entityID *string
	if rawEID := c.Query("entity_id"); rawEID != "" {
		entityID = &rawEID
	}

	var startDate *time.Time
	if rawStart := c.Query("start_date"); rawStart != "" {
		if t, err := time.Parse("2006-01-02", rawStart); err == nil {
			startDate = &t
		}
	}

	var endDate *time.Time
	if rawEnd := c.Query("end_date"); rawEnd != "" {
		if t, err := time.Parse("2006-01-02", rawEnd); err == nil {
			// Set end of day
			tEnd := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &tEnd
		}
	}

	filter := domain.AuditLogFilter{
		UserID:    userID,
		Action:    c.Query("action"),
		Entity:    c.Query("entity"),
		EntityID:  entityID,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PerPage:   perPage,
	}

	logs, meta, err := h.auditUsecase.ListAuditLogs(c.Request.Context(), filter)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat jejak audit: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Jejak audit log berhasil dimuat", logs, meta)
}
