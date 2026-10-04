package http

import (
	"net/http"
	"strconv"

	"backend-alusi-go/internal/delivery/http/middleware"
	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	catalogUsecase *usecase.CatalogUsecase
}

func NewCatalogHandler(catalogUsecase *usecase.CatalogUsecase) *CatalogHandler {
	return &CatalogHandler{
		catalogUsecase: catalogUsecase,
	}
}

// ListCategories godoc
// @Summary      Daftar kategori aplikasi
// @Description  Mengambil seluruh kategori aplikasi yang aktif
// @Tags         Catalog
// @Produce      json
// @Success      200  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /categories [get]
func (h *CatalogHandler) ListCategories(c *gin.Context) {
	categories, err := h.catalogUsecase.ListCategories(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal memuat kategori aplikasi: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Daftar kategori berhasil dimuat", categories, nil)
}

// GetCategoryBySlug godoc
// @Summary      Detail kategori berdasarkan slug
// @Description  Mengambil detail data sebuah kategori aplikasi
// @Tags         Catalog
// @Produce      json
// @Param        slug  path      string  true  "Slug kategori"
// @Success      200   {object}  response.StandardResponse
// @Failure      404   {object}  response.StandardResponse
// @Router       /categories/{slug} [get]
func (h *CatalogHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	category, err := h.catalogUsecase.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil || category == nil {
		response.NotFound(c, "Kategori aplikasi tidak ditemukan.")
		return
	}

	response.Success(c, http.StatusOK, "Detail kategori berhasil dimuat", category, nil)
}

// ListApps godoc
// @Summary      Daftar katalog aplikasi
// @Description  Mengambil daftar aplikasi dengan filter kategori, target pengguna, status, persona, dan pagination
// @Tags         Catalog
// @Produce      json
// @Param        category         query     string  false  "Filter slug kategori"
// @Param        target_pengguna  query     string  false  "Filter target pengguna (semua, internal, internal_bps, mitra, pimpinan)"
// @Param        status           query     string  false  "Filter status operasional (operasional, pemeliharaan, kendala)"
// @Param        persona          query     string  false  "Filter persona / tag"
// @Param        page             query     int     false  "Halaman (default 1)"
// @Param        per_page         query     int     false  "Jumlah item per halaman (default 20)"
// @Success      200              {object}  response.StandardResponse
// @Failure      500              {object}  response.StandardResponse
// @Router       /apps [get]
func (h *CatalogHandler) ListApps(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	var claims *jwt.SessionClaims
	if val, exists := c.Get(middleware.CtxUserClaims); exists {
		claims = val.(*jwt.SessionClaims)
	}

	filter := usecase.CatalogFilterInput{
		CategorySlug:   c.Query("category"),
		TargetPengguna: c.Query("target_pengguna"),
		StatusLayanan:  c.Query("status"),
		Persona:        c.Query("persona"),
		Page:           page,
		PerPage:        perPage,
	}

	apps, meta, err := h.catalogUsecase.ListApps(c.Request.Context(), filter, claims)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat katalog aplikasi: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Katalog aplikasi berhasil dimuat", apps, meta)
}

// GetAppBySlug godoc
// @Summary      Detail aplikasi berdasarkan slug
// @Description  Mengambil metadata lengkap, deskripsi, panduan, dan status sebuah aplikasi
// @Tags         Catalog
// @Produce      json
// @Param        slug  path      string  true  "Slug aplikasi"
// @Success      200   {object}  response.StandardResponse
// @Failure      404   {object}  response.StandardResponse
// @Router       /apps/{slug} [get]
func (h *CatalogHandler) GetAppBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var userID *int
	if val, exists := c.Get(middleware.CtxUserID); exists {
		id := val.(int)
		userID = &id
	}

	app, err := h.catalogUsecase.GetAppBySlug(c.Request.Context(), slug, userID)
	if err != nil || app == nil {
		response.NotFound(c, "Aplikasi tidak ditemukan.")
		return
	}

	response.Success(c, http.StatusOK, "Detail aplikasi berhasil dimuat", app, nil)
}

// SearchApps godoc
// @Summary      Pencarian cerdas aplikasi (Trigram pg_trgm)
// @Description  Pencarian aplikasi dengan toleransi typo, sinonim kata kunci, dan filter kategori
// @Tags         Catalog
// @Produce      json
// @Param        q         query     string  true   "Kata kunci pencarian"
// @Param        category  query     string  false  "Filter slug kategori"
// @Param        limit     query     int     false  "Batas jumlah hasil (default 20)"
// @Success      200       {object}  response.StandardResponse
// @Failure      400       {object}  response.StandardResponse
// @Failure      500       {object}  response.StandardResponse
// @Router       /apps/search [get]
func (h *CatalogHandler) SearchApps(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Kata kunci pencarian (q) tidak boleh kosong.", nil)
		return
	}

	categorySlug := c.Query("category")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var claims *jwt.SessionClaims
	if val, exists := c.Get(middleware.CtxUserClaims); exists {
		claims = val.(*jwt.SessionClaims)
	}

	apps, err := h.catalogUsecase.SearchApps(c.Request.Context(), query, categorySlug, claims, limit)
	if err != nil {
		response.InternalServerError(c, "Gagal mencari aplikasi: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Hasil pencarian aplikasi berhasil dimuat", apps, nil)
}

// GetAppGuides godoc
// @Summary      Panduan dan FAQ aplikasi
// @Description  Mengambil daftar petunjuk teknis / user manual / FAQ untuk aplikasi terkait
// @Tags         Catalog
// @Produce      json
// @Param        slug  path      string  true  "Slug aplikasi"
// @Success      200   {object}  response.StandardResponse
// @Failure      404   {object}  response.StandardResponse
// @Router       /apps/{slug}/guides [get]
func (h *CatalogHandler) GetAppGuides(c *gin.Context) {
	slug := c.Param("slug")
	guides, err := h.catalogUsecase.GetGuidesByAppSlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Panduan aplikasi berhasil dimuat", guides, nil)
}
