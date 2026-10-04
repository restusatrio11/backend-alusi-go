package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
)

type AnnouncementUsecase struct {
	announcementRepo *postgres.AnnouncementRepo
	appRepo          *postgres.AppRepo
}

func NewAnnouncementUsecase(
	announcementRepo *postgres.AnnouncementRepo,
	appRepo *postgres.AppRepo,
) *AnnouncementUsecase {
	return &AnnouncementUsecase{
		announcementRepo: announcementRepo,
		appRepo:          appRepo,
	}
}

func (u *AnnouncementUsecase) GetActiveAnnouncements(ctx context.Context, appID *int) ([]domain.Announcement, error) {
	if u.announcementRepo == nil {
		return []domain.Announcement{}, nil
	}
	return u.announcementRepo.ListActive(ctx, appID)
}

func (u *AnnouncementUsecase) ListAdminAnnouncements(ctx context.Context, page, perPage int) ([]domain.Announcement, *response.Meta, error) {
	if u.announcementRepo == nil {
		return []domain.Announcement{}, &response.Meta{Page: 1, PerPage: perPage, TotalItems: 0, TotalPages: 0}, nil
	}

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	list, total, err := u.announcementRepo.ListAll(ctx, perPage, offset)
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

type CreateAnnouncementInput struct {
	Judul    string     `json:"judul" binding:"required"`
	Konten   string     `json:"konten" binding:"required"`
	Tipe     string     `json:"tipe"` // info, warning, maintenance, urgent
	AppID    *int       `json:"app_id"`
	IsActive bool       `json:"is_active"`
	StartsAt *time.Time `json:"starts_at"`
	EndsAt   *time.Time `json:"ends_at"`
}

type UpdateAnnouncementInput struct {
	Judul    string     `json:"judul" binding:"required"`
	Konten   string     `json:"konten" binding:"required"`
	Tipe     string     `json:"tipe"`
	AppID    *int       `json:"app_id"`
	IsActive bool       `json:"is_active"`
	StartsAt *time.Time `json:"starts_at"`
	EndsAt   *time.Time `json:"ends_at"`
}

func (u *AnnouncementUsecase) CreateAnnouncement(ctx context.Context, input CreateAnnouncementInput) (*domain.Announcement, error) {
	tipe := input.Tipe
	if tipe == "" {
		tipe = "info"
	}

	startsAt := time.Now()
	if input.StartsAt != nil && !input.StartsAt.IsZero() {
		startsAt = *input.StartsAt
	}

	ann := &domain.Announcement{
		Judul:    input.Judul,
		Konten:   input.Konten,
		Tipe:     tipe,
		AppID:    input.AppID,
		IsActive: input.IsActive,
		StartsAt: startsAt,
		EndsAt:   input.EndsAt,
	}

	if err := u.announcementRepo.Create(ctx, ann); err != nil {
		return nil, fmt.Errorf("gagal membuat pengumuman: %w", err)
	}

	return ann, nil
}

func (u *AnnouncementUsecase) UpdateAnnouncement(ctx context.Context, id int, input UpdateAnnouncementInput) (*domain.Announcement, error) {
	ann, err := u.announcementRepo.GetByID(ctx, id)
	if err != nil || ann == nil {
		return nil, fmt.Errorf("pengumuman tidak ditemukan")
	}

	tipe := input.Tipe
	if tipe == "" {
		tipe = "info"
	}

	startsAt := ann.StartsAt
	if input.StartsAt != nil && !input.StartsAt.IsZero() {
		startsAt = *input.StartsAt
	}

	ann.Judul = input.Judul
	ann.Konten = input.Konten
	ann.Tipe = tipe
	ann.AppID = input.AppID
	ann.IsActive = input.IsActive
	ann.StartsAt = startsAt
	ann.EndsAt = input.EndsAt

	if err := u.announcementRepo.Update(ctx, ann); err != nil {
		return nil, fmt.Errorf("gagal memperbarui pengumuman: %w", err)
	}

	return ann, nil
}

func (u *AnnouncementUsecase) DeleteAnnouncement(ctx context.Context, id int) error {
	return u.announcementRepo.Delete(ctx, id)
}
