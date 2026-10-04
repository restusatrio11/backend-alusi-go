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

type InteractionHandler struct {
	interactionUsecase *usecase.InteractionUsecase
}

func NewInteractionHandler(interactionUsecase *usecase.InteractionUsecase) *InteractionHandler {
	return &InteractionHandler{
		interactionUsecase: interactionUsecase,
	}
}

// RecordClick handles click tracking for opening an application
func (h *InteractionHandler) RecordClick(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID aplikasi tidak valid.", nil)
		return
	}

	var claims *jwt.SessionClaims
	if val, exists := c.Get(middleware.CtxUserClaims); exists {
		claims = val.(*jwt.SessionClaims)
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	app, err := h.interactionUsecase.RecordClick(c.Request.Context(), appID, claims, ipAddress, userAgent)
	if err != nil {
		response.NotFound(c, "Aplikasi tidak ditemukan.")
		return
	}

	response.Success(c, http.StatusOK, "Klik berhasil dicatat", gin.H{
		"app_id": app.ID,
		"nama":   app.Nama,
		"url":    app.URL,
	}, nil)
}

// ToggleFavorite adds or removes an app from user favorites
func (h *InteractionHandler) ToggleFavorite(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID aplikasi tidak valid.", nil)
		return
	}

	userID := c.GetInt(middleware.CtxUserID)
	if userID == 0 {
		response.Unauthorized(c, "Silakan login terlebih dahulu untuk menambah favorit.")
		return
	}

	isFavorite, err := h.interactionUsecase.ToggleFavorite(c.Request.Context(), userID, appID)
	if err != nil {
		response.InternalServerError(c, "Gagal mengubah status favorit: "+err.Error())
		return
	}

	message := "Aplikasi berhasil ditambahkan ke favorit"
	if !isFavorite {
		message = "Aplikasi berhasil dihapus dari favorit"
	}

	response.Success(c, http.StatusOK, message, gin.H{
		"app_id":      appID,
		"is_favorite": isFavorite,
	}, nil)
}

// ListFavorites returns all favorite apps of the logged in user
func (h *InteractionHandler) ListFavorites(c *gin.Context) {
	userID := c.GetInt(middleware.CtxUserID)
	if userID == 0 {
		response.Unauthorized(c, "Pengguna belum terautentikasi.")
		return
	}

	apps, err := h.interactionUsecase.ListFavorites(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat daftar favorit: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Daftar aplikasi favorit berhasil dimuat", apps, nil)
}

// ListRecents returns recently clicked apps of the logged in user
func (h *InteractionHandler) ListRecents(c *gin.Context) {
	userID := c.GetInt(middleware.CtxUserID)
	if userID == 0 {
		response.Unauthorized(c, "Pengguna belum terautentikasi.")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	apps, err := h.interactionUsecase.ListRecents(c.Request.Context(), userID, limit)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat riwayat aplikasi terakhir dibuka: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Riwayat aplikasi terakhir dibuka berhasil dimuat", apps, nil)
}
