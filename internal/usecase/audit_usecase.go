package usecase

import (
	"context"
	"math"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
)

type AuditUsecase struct {
	auditRepo *postgres.AuditLogRepo
}

func NewAuditUsecase(auditRepo *postgres.AuditLogRepo) *AuditUsecase {
	return &AuditUsecase{
		auditRepo: auditRepo,
	}
}

func (u *AuditUsecase) Record(ctx context.Context, entry *domain.AuditLog) {
	if u.auditRepo == nil || entry == nil {
		return
	}
	// Asynchronous/non-blocking or direct recording
	_ = u.auditRepo.Record(ctx, entry)
}

func (u *AuditUsecase) ListAuditLogs(ctx context.Context, filter domain.AuditLogFilter) ([]domain.AuditLog, *response.Meta, error) {
	if u.auditRepo == nil {
		return []domain.AuditLog{}, &response.Meta{Page: 1, PerPage: filter.PerPage, TotalItems: 0, TotalPages: 0}, nil
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}
	filter.Page = page
	filter.PerPage = perPage

	logs, total, err := u.auditRepo.List(ctx, filter)
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

	return logs, meta, nil
}
