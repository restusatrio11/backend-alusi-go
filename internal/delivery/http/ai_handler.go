package http

import (
	"net/http"

	"backend-alusi-go/internal/delivery/http/middleware"
	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiUsecase *usecase.AIUsecase
}

func NewAIHandler(aiUsecase *usecase.AIUsecase) *AIHandler {
	return &AIHandler{
		aiUsecase: aiUsecase,
	}
}

type AskRequest struct {
	Question string `json:"question" binding:"required"`
}

// Ask handles natural language questions and recommends relevant apps
func (h *AIHandler) Ask(c *gin.Context) {
	var req AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Pertanyaan (question) tidak boleh kosong.", nil)
		return
	}

	var claims *jwt.SessionClaims
	if val, exists := c.Get(middleware.CtxUserClaims); exists {
		claims = val.(*jwt.SessionClaims)
	}

	res, err := h.aiUsecase.Ask(c.Request.Context(), req.Question, claims)
	if err != nil {
		response.InternalServerError(c, "Gagal memproses rekomendasi AI: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Rekomendasi asisten AI berhasil ditemukan", res, nil)
}
