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

// ListCategories returns all application categories
func (h *CatalogHandler) ListCategories(c *gin.Context) {
	categories, err := h.catalogUsecase.ListCategories(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal memuat kategori aplikasi: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Daftar kategori berhasil dimuat", categories, nil)
}

// GetCategoryBySlug returns category details by slug
func (h *CatalogHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	category, err := h.catalogUsecase.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil || category == nil {
		response.NotFound(c, "Kategori aplikasi tidak ditemukan.")
		return
	}

	response.Success(c, http.StatusOK, "Detail kategori berhasil dimuat", category, nil)
}

// ListApps returns application catalog with filtering and pagination
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

// GetAppBySlug returns detailed application metadata by slug
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

// SearchApps handles full text and trigram search with typo tolerance
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

// GetAppGuides returns guides and FAQ documentation for an application
func (h *CatalogHandler) GetAppGuides(c *gin.Context) {
	slug := c.Param("slug")
	guides, err := h.catalogUsecase.GetGuidesByAppSlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Panduan aplikasi berhasil dimuat", guides, nil)
}
