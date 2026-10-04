package http

import (
	"net/http"
	"strconv"

	"backend-alusi-go/internal/delivery/http/middleware"
	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type FeedbackHandler struct {
	feedbackUsecase *usecase.FeedbackUsecase
}

func NewFeedbackHandler(feedbackUsecase *usecase.FeedbackUsecase) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackUsecase: feedbackUsecase,
	}
}

// SubmitFeedback handles submission of user feedback or issue report
func (h *FeedbackHandler) SubmitFeedback(c *gin.Context) {
	var input usecase.SubmitFeedbackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload laporan tidak valid: "+err.Error(), nil)
		return
	}

	var userID *int
	if val, exists := c.Get(middleware.CtxUserID); exists {
		id := val.(int)
		userID = &id
	}

	fb, err := h.feedbackUsecase.SubmitFeedback(c.Request.Context(), input, userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Laporan/saran Anda berhasil dikirim dan akan ditindaklanjuti.", fb, nil)
}

// ListAdminFeedbacks handles listing feedbacks for admin
func (h *FeedbackHandler) ListAdminFeedbacks(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	var appID *int
	if rawAppID := c.Query("app_id"); rawAppID != "" {
		if id, err := strconv.Atoi(rawAppID); err == nil {
			appID = &id
		}
	}

	list, meta, err := h.feedbackUsecase.ListFeedbacks(c.Request.Context(), status, appID, page, perPage)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat daftar feedback: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Daftar laporan feedback berhasil dimuat", list, meta)
}

// UpdateFeedbackStatus handles resolving/updating feedback status
func (h *FeedbackHandler) UpdateFeedbackStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID laporan tidak valid.", nil)
		return
	}

	var input usecase.UpdateFeedbackStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload update feedback tidak valid: "+err.Error(), nil)
		return
	}

	fb, err := h.feedbackUsecase.UpdateFeedbackStatus(c.Request.Context(), id, input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Status laporan feedback berhasil diperbarui", fb, nil)
}
