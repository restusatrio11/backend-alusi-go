package http

import (
	"net/http"
	"strconv"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/media"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUsecase   *usecase.AdminUsecase
	imageOptimizer *media.ImageOptimizer
}

func NewAdminHandler(adminUsecase *usecase.AdminUsecase, imageOptimizer *media.ImageOptimizer) *AdminHandler {
	return &AdminHandler{
		adminUsecase:   adminUsecase,
		imageOptimizer: imageOptimizer,
	}
}

// CreateApp godoc
// @Summary      Tambah aplikasi baru
// @Description  Menambahkan data aplikasi baru ke dalam katalog (Admin)
// @Tags         Admin - Catalog
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      usecase.CreateAppInput  true  "Data aplikasi baru"
// @Success      201    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/apps [post]
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

// UpdateApp godoc
// @Summary      Perbarui data aplikasi
// @Description  Memperbarui informasi metadata aplikasi yang sudah ada (Admin)
// @Tags         Admin - Catalog
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                     true  "ID Aplikasi"
// @Param        input  body      usecase.UpdateAppInput  true  "Data perubahan aplikasi"
// @Success      200    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/apps/{id} [put]
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

// DeleteApp godoc
// @Summary      Hapus (Soft Delete) aplikasi
// @Description  Menonaktifkan aplikasi dari tampilan katalog pengguna (Admin)
// @Tags         Admin - Catalog
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID Aplikasi"
// @Success      200  {object}  response.StandardResponse
// @Failure      400  {object}  response.StandardResponse
// @Failure      401  {object}  response.StandardResponse
// @Failure      403  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /admin/apps/{id} [delete]
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

// ReorderApps godoc
// @Summary      Ubah urutan (Drag & Drop) aplikasi
// @Description  Memperbarui urutan sequence posisi tampilan katalog aplikasi (Admin)
// @Tags         Admin - Catalog
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      ReorderRequest  true  "Daftar ID aplikasi sesuai urutan baru"
// @Success      200    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/apps/reorder [put]
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

// CreateCategory godoc
// @Summary      Tambah kategori baru
// @Description  Menambahkan kategori kelompok aplikasi baru (Admin)
// @Tags         Admin - Catalog
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      usecase.CategoryInput  true  "Data kategori baru"
// @Success      201    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/categories [post]
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

// UpdateCategory godoc
// @Summary      Perbarui kategori
// @Description  Memperbarui nama, deskripsi, ikon, atau warna kategori (Admin)
// @Tags         Admin - Catalog
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                    true  "ID Kategori"
// @Param        input  body      usecase.CategoryInput  true  "Data perubahan kategori"
// @Success      200    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/categories/{id} [put]
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

// DeleteCategory godoc
// @Summary      Hapus kategori
// @Description  Menghapus kategori aplikasi (Admin)
// @Tags         Admin - Catalog
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID Kategori"
// @Success      200  {object}  response.StandardResponse
// @Failure      400  {object}  response.StandardResponse
// @Failure      401  {object}  response.StandardResponse
// @Failure      403  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /admin/categories/{id} [delete]
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

// CreateGuide godoc
// @Summary      Tambah panduan / FAQ aplikasi
// @Description  Menambahkan panduan atau FAQ baru untuk sebuah aplikasi (Admin)
// @Tags         Admin - Guides
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                 true  "ID Aplikasi"
// @Param        input  body      usecase.GuideInput  true  "Data panduan baru"
// @Success      201    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/apps/{id}/guides [post]
func (h *AdminHandler) CreateGuide(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID aplikasi tidak valid.", nil)
		return
	}

	var input usecase.GuideInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload panduan tidak valid: "+err.Error(), nil)
		return
	}

	guide, err := h.adminUsecase.CreateGuide(c.Request.Context(), appID, input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Panduan aplikasi berhasil ditambahkan", guide, nil)
}

// UpdateGuide godoc
// @Summary      Perbarui panduan / FAQ
// @Description  Memperbarui judul, konten markdown, tipe, atau urutan panduan (Admin)
// @Tags         Admin - Guides
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                 true  "ID Panduan"
// @Param        input  body      usecase.GuideInput  true  "Data perubahan panduan"
// @Success      200    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/guides/{id} [put]
func (h *AdminHandler) UpdateGuide(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID panduan tidak valid.", nil)
		return
	}

	var input usecase.GuideInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload panduan tidak valid: "+err.Error(), nil)
		return
	}

	guide, err := h.adminUsecase.UpdateGuide(c.Request.Context(), id, input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Panduan aplikasi berhasil diperbarui", guide, nil)
}

// DeleteGuide godoc
// @Summary      Hapus panduan / FAQ
// @Description  Menghapus data panduan atau FAQ (Admin)
// @Tags         Admin - Guides
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID Panduan"
// @Success      200  {object}  response.StandardResponse
// @Failure      400  {object}  response.StandardResponse
// @Failure      401  {object}  response.StandardResponse
// @Failure      403  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /admin/guides/{id} [delete]
func (h *AdminHandler) DeleteGuide(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID panduan tidak valid.", nil)
		return
	}

	if err := h.adminUsecase.DeleteGuide(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, "Gagal menghapus panduan: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Panduan aplikasi berhasil dihapus", nil, nil)
}

// UploadAppLogo godoc
// @Summary      Upload dan kompresi logo aplikasi
// @Description  Mengunggah file gambar logo aplikasi (PNG/JPG/JPEG/WebP/SVG, maks 2MB), otomatis di-resize maks 512x512 dan dikompresi (Admin)
// @Tags         Admin - Catalog
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int   true  "ID Aplikasi"
// @Param        logo  formData  file  true  "File logo gambar (PNG, JPG, JPEG, WebP, SVG maks 2MB)"
// @Success      200   {object}  response.StandardResponse
// @Failure      400   {object}  response.StandardResponse
// @Failure      401   {object}  response.StandardResponse
// @Failure      403   {object}  response.StandardResponse
// @Failure      404   {object}  response.StandardResponse
// @Failure      500   {object}  response.StandardResponse
// @Router       /admin/apps/{id}/logo [post]
func (h *AdminHandler) UploadAppLogo(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID aplikasi tidak valid.", nil)
		return
	}

	app, err := h.adminUsecase.GetAppByID(c.Request.Context(), appID)
	if err != nil || app == nil {
		response.NotFound(c, "Aplikasi tidak ditemukan.")
		return
	}

	if h.imageOptimizer == nil {
		response.InternalServerError(c, "Image optimizer service belum dikonfigurasi.")
		return
	}

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		fileHeader, err = c.FormFile("file")
	}
	if err != nil {
		fileHeader, err = c.FormFile("image")
	}
	if err != nil || fileHeader == nil {
		response.BadRequest(c, "File logo tidak ditemukan dalam request form (field 'logo')", nil)
		return
	}

	logoURL, err := h.imageOptimizer.ProcessAndSave(fileHeader, app.Slug)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	updatedApp, err := h.adminUsecase.UpdateAppLogo(c.Request.Context(), appID, logoURL)
	if err != nil {
		response.InternalServerError(c, "Gagal memperbarui URL logo aplikasi: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Logo aplikasi berhasil diunggah dan dikompresi", map[string]interface{}{
		"app":      updatedApp,
		"logo_url": logoURL,
	}, nil)
}

