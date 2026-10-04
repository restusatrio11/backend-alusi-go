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

// GetActiveAnnouncements godoc
// @Summary      Daftar pengumuman / banner aktif
// @Description  Mengambil pengumuman broadcast sistem atau banner pemeliharaan spesifik aplikasi
// @Tags         Announcements
// @Produce      json
// @Param        app_id  query     int  false  "Filter spesifik ID aplikasi (opsional)"
// @Success      200     {object}  response.StandardResponse
// @Failure      500     {object}  response.StandardResponse
// @Router       /announcements [get]
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

// ListAdminAnnouncements godoc
// @Summary      Daftar seluruh pengumuman (Admin)
// @Description  Mengambil seluruh riwayat pengumuman baik yang aktif maupun nonaktif (Admin)
// @Tags         Admin - Announcements
// @Produce      json
// @Security     BearerAuth
// @Param        page      query     int  false  "Halaman (default 1)"
// @Param        per_page  query     int  false  "Jumlah item per halaman (default 20)"
// @Success      200       {object}  response.StandardResponse
// @Failure      401       {object}  response.StandardResponse
// @Failure      403       {object}  response.StandardResponse
// @Failure      500       {object}  response.StandardResponse
// @Router       /admin/announcements [get]
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

// CreateAnnouncement godoc
// @Summary      Buat pengumuman baru
// @Description  Membuat pengumuman broadcast atau banner pemeliharaan baru (Admin)
// @Tags         Admin - Announcements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      usecase.CreateAnnouncementInput  true  "Data pengumuman baru"
// @Success      201    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/announcements [post]
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

// UpdateAnnouncement godoc
// @Summary      Perbarui pengumuman
// @Description  Memperbarui judul, pesan, tipe, atau masa tayang pengumuman (Admin)
// @Tags         Admin - Announcements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                              true  "ID Pengumuman"
// @Param        input  body      usecase.UpdateAnnouncementInput  true  "Data perubahan pengumuman"
// @Success      200    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/announcements/{id} [put]
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

// DeleteAnnouncement godoc
// @Summary      Hapus pengumuman
// @Description  Menghapus pengumuman dari sistem (Admin)
// @Tags         Admin - Announcements
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID Pengumuman"
// @Success      200  {object}  response.StandardResponse
// @Failure      400  {object}  response.StandardResponse
// @Failure      401  {object}  response.StandardResponse
// @Failure      403  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /admin/announcements/{id} [delete]
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
