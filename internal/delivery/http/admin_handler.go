package http

import (
	"net/http"
	"strconv"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUsecase *usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase *usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		adminUsecase: adminUsecase,
	}
}

// CreateApp handles application creation
func (h *AdminHandler) CreateApp(c *gin.Context) {
	var input usecase.CreateAppInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload aplikasi tidak valid: "+err.Error(), nil)
		return
	}

	app, err := h.adminUsecase.CreateApp(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Aplikasi berhasil ditambahkan", app, nil)
}

// UpdateApp handles application update
func (h *AdminHandler) UpdateApp(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID aplikasi tidak valid.", nil)
		return
	}

	var input usecase.UpdateAppInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload aplikasi tidak valid: "+err.Error(), nil)
		return
	}

	app, err := h.adminUsecase.UpdateApp(c.Request.Context(), id, input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Aplikasi berhasil diperbarui", app, nil)
}

// DeleteApp handles application soft deletion
func (h *AdminHandler) DeleteApp(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID aplikasi tidak valid.", nil)
		return
	}

	if err := h.adminUsecase.DeleteApp(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, "Gagal menonaktifkan aplikasi: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Aplikasi berhasil dinonaktifkan", nil, nil)
}

type ReorderRequest struct {
	AppIDs []int `json:"app_ids" binding:"required"`
}

// ReorderApps handles drag and drop reordering of applications
func (h *AdminHandler) ReorderApps(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Daftar urutan ID aplikasi (app_ids) harus berupa array angka.", nil)
		return
	}

	if err := h.adminUsecase.ReorderApps(c.Request.Context(), req.AppIDs); err != nil {
		response.InternalServerError(c, "Gagal mengurutkan aplikasi: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Urutan aplikasi berhasil diperbarui", nil, nil)
}

// CreateCategory handles category creation
func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var input usecase.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload kategori tidak valid: "+err.Error(), nil)
		return
	}

	cat, err := h.adminUsecase.CreateCategory(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Kategori berhasil dibuat", cat, nil)
}

// UpdateCategory handles category update
func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID kategori tidak valid.", nil)
		return
	}

	var input usecase.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload kategori tidak valid: "+err.Error(), nil)
		return
	}

	cat, err := h.adminUsecase.UpdateCategory(c.Request.Context(), id, input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Kategori berhasil diperbarui", cat, nil)
}

// DeleteCategory handles category deletion
func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID kategori tidak valid.", nil)
		return
	}

	if err := h.adminUsecase.DeleteCategory(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, "Gagal menghapus kategori: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Kategori berhasil dihapus", nil, nil)
}
