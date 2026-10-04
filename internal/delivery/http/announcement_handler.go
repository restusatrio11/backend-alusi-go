package http

import (
	"net/http"
	"strconv"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	announcementUsecase *usecase.AnnouncementUsecase
}

func NewAnnouncementHandler(announcementUsecase *usecase.AnnouncementUsecase) *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementUsecase: announcementUsecase,
	}
}

// GetActiveAnnouncements returns currently active announcements (Public)
func (h *AnnouncementHandler) GetActiveAnnouncements(c *gin.Context) {
	var appID *int
	if rawAppID := c.Query("app_id"); rawAppID != "" {
		if id, err := strconv.Atoi(rawAppID); err == nil {
			appID = &id
		}
	}

	list, err := h.announcementUsecase.GetActiveAnnouncements(c.Request.Context(), appID)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat pengumuman: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Daftar pengumuman aktif berhasil dimuat", list, nil)
}

// ListAdminAnnouncements returns all announcements for admin management
func (h *AnnouncementHandler) ListAdminAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	list, meta, err := h.announcementUsecase.ListAdminAnnouncements(c.Request.Context(), page, perPage)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat daftar pengumuman: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Daftar pengumuman admin berhasil dimuat", list, meta)
}

// CreateAnnouncement handles announcement creation by admin
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	var input usecase.CreateAnnouncementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload pengumuman tidak valid: "+err.Error(), nil)
		return
	}

	ann, err := h.announcementUsecase.CreateAnnouncement(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Pengumuman berhasil dibuat", ann, nil)
}

// UpdateAnnouncement handles announcement update by admin
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID pengumuman tidak valid.", nil)
		return
	}

	var input usecase.UpdateAnnouncementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload pengumuman tidak valid: "+err.Error(), nil)
		return
	}

	ann, err := h.announcementUsecase.UpdateAnnouncement(c.Request.Context(), id, input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Pengumuman berhasil diperbarui", ann, nil)
}

// DeleteAnnouncement handles announcement deletion by admin
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID pengumuman tidak valid.", nil)
		return
	}

	if err := h.announcementUsecase.DeleteAnnouncement(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, "Gagal menghapus pengumuman: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Pengumuman berhasil dihapus", nil, nil)
}
