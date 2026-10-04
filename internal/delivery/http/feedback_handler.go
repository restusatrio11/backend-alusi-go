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

// SubmitFeedback godoc
// @Summary      Kirim kritik, saran, atau laporan kendala
// @Description  Mengirim formulir laporan kendala atau usulan fitur untuk aplikasi tertentu atau portal secara umum
// @Tags         Feedback
// @Accept       json
// @Produce      json
// @Param        input  body      usecase.SubmitFeedbackInput  true  "Data formulir feedback"
// @Success      201    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /feedbacks [post]
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

// ListAdminFeedbacks godoc
// @Summary      Daftar laporan feedback (Admin)
// @Description  Mengambil seluruh tiket laporan kendala, saran, dan pertanyaan pengguna (Admin)
// @Tags         Admin - Feedback
// @Produce      json
// @Security     BearerAuth
// @Param        status    query     string  false  "Filter status tiket (pending, in_progress, resolved, closed)"
// @Param        app_id    query     int     false  "Filter ID aplikasi"
// @Param        page      query     int     false  "Halaman (default 1)"
// @Param        per_page  query     int     false  "Jumlah per halaman (default 20)"
// @Success      200       {object}  response.StandardResponse
// @Failure      401       {object}  response.StandardResponse
// @Failure      403       {object}  response.StandardResponse
// @Failure      500       {object}  response.StandardResponse
// @Router       /admin/feedbacks [get]
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

// UpdateFeedbackStatus godoc
// @Summary      Perbarui status tindak lanjut tiket feedback (Admin)
// @Description  Mengubah status dan menambahkan catatan respon tindak lanjut pada tiket laporan (Admin)
// @Tags         Admin - Feedback
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                                 true  "ID Feedback"
// @Param        input  body      usecase.UpdateFeedbackStatusInput  true  "Data status dan respon tanggapan"
// @Success      200    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      401    {object}  response.StandardResponse
// @Failure      403    {object}  response.StandardResponse
// @Failure      500    {object}  response.StandardResponse
// @Router       /admin/feedbacks/{id} [put]
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
