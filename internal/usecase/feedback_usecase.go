package usecase

import (
	"context"
	"fmt"
	"math"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
)

type FeedbackUsecase struct {
	feedbackRepo *postgres.FeedbackRepo
	appRepo      *postgres.AppRepo
}

func NewFeedbackUsecase(
	feedbackRepo *postgres.FeedbackRepo,
	appRepo *postgres.AppRepo,
) *FeedbackUsecase {
	return &FeedbackUsecase{
		feedbackRepo: feedbackRepo,
		appRepo:      appRepo,
	}
}

type SubmitFeedbackInput struct {
	AppID    *int   `json:"app_id"`
	Kategori string `json:"kategori"` // kendala, saran, data
	Pesan    string `json:"pesan" binding:"required"`
}

func (u *FeedbackUsecase) SubmitFeedback(ctx context.Context, input SubmitFeedbackInput, userID *int) (*domain.Feedback, error) {
	kategori := input.Kategori
	if kategori == "" {
		kategori = "kendala"
	}

	fb := &domain.Feedback{
		UserID:   userID,
		AppID:    input.AppID,
		Kategori: kategori,
		Pesan:    input.Pesan,
		Status:   "pending",
	}

	if u.feedbackRepo != nil {
		if err := u.feedbackRepo.Create(ctx, fb); err != nil {
			return nil, fmt.Errorf("gagal mengirim laporan: %w", err)
		}
	}

	return fb, nil
}

func (u *FeedbackUsecase) ListFeedbacks(ctx context.Context, status string, appID *int, page, perPage int) ([]domain.Feedback, *response.Meta, error) {
	if u.feedbackRepo == nil {
		return []domain.Feedback{}, &response.Meta{Page: 1, PerPage: perPage, TotalItems: 0, TotalPages: 0}, nil
	}

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	list, total, err := u.feedbackRepo.List(ctx, status, appID, perPage, offset)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	meta := &response.Meta{
		Page:       page,
		PerPage:    perPage,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return list, meta, nil
}

type UpdateFeedbackStatusInput struct {
	Status       string  `json:"status" binding:"required"` // pending, in_progress, resolved, closed
	CatatanAdmin *string `json:"catatan_admin"`
}

func (u *FeedbackUsecase) UpdateFeedbackStatus(ctx context.Context, id int, input UpdateFeedbackStatusInput) (*domain.Feedback, error) {
	if u.feedbackRepo == nil {
		return nil, fmt.Errorf("repository feedback tidak tersedia")
	}

	fb, err := u.feedbackRepo.GetByID(ctx, id)
	if err != nil || fb == nil {
		return nil, fmt.Errorf("laporan feedback tidak ditemukan")
	}

	if err := u.feedbackRepo.UpdateStatus(ctx, id, input.Status, input.CatatanAdmin); err != nil {
		return nil, fmt.Errorf("gagal memperbarui status feedback: %w", err)
	}

	fb.Status = input.Status
	fb.CatatanAdmin = input.CatatanAdmin

	return fb, nil
}
