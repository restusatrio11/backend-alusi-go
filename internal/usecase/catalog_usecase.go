package usecase

import (
	"context"
	"fmt"
	"math"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
	"backend-alusi-go/pkg/jwt"
)

type CatalogUsecase struct {
	categoryRepo *postgres.CategoryRepo
	appRepo      *postgres.AppRepo
	userRepo     *postgres.UserRepo
}

func NewCatalogUsecase(
	categoryRepo *postgres.CategoryRepo,
	appRepo      *postgres.AppRepo,
	userRepo     *postgres.UserRepo,
) *CatalogUsecase {
	return &CatalogUsecase{
		categoryRepo: categoryRepo,
		appRepo:      appRepo,
		userRepo:     userRepo,
	}
}

func (u *CatalogUsecase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return u.categoryRepo.List(ctx)
}

func (u *CatalogUsecase) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	return u.categoryRepo.GetBySlug(ctx, slug)
}

type CatalogFilterInput struct {
	CategorySlug   string
	TargetPengguna string
	StatusLayanan  string
	Persona        string // "lapangan", "pegawai", "pimpinan", "mitra"
	Page           int
	PerPage        int
}

func (u *CatalogUsecase) ListApps(
	ctx context.Context,
	input CatalogFilterInput,
	claims *jwt.SessionClaims,
) ([]domain.App, *response.Meta, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}

	perPage := input.PerPage
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	var categoryID *int
	if input.CategorySlug != "" {
		cat, err := u.categoryRepo.GetBySlug(ctx, input.CategorySlug)
		if err == nil && cat != nil {
			categoryID = &cat.ID
		}
	}

	var targetPengguna *string
	if input.TargetPengguna != "" {
		targetPengguna = &input.TargetPengguna
	} else if input.Persona != "" {
		switch input.Persona {
		case "lapangan":
			p := "Petugas Lapangan"
			targetPengguna = &p
		case "pegawai":
			p := "Pegawai"
			targetPengguna = &p
		case "pimpinan":
			p := "Pimpinan"
			targetPengguna = &p
		case "mitra":
			p := "Mitra"
			targetPengguna = &p
		}
	}

	var statusLayanan *string
	if input.StatusLayanan != "" {
		statusLayanan = &input.StatusLayanan
	}

	var userID *int
	var roleIDs []int
	if claims != nil {
		userID = &claims.UserID
		// Resolve user's database role IDs for accurate app_access filtering
		user, err := u.userRepo.GetByID(ctx, claims.UserID)
		if err == nil && user != nil {
			for _, r := range user.Roles {
				roleIDs = append(roleIDs, r.ID)
			}
		}
	}

	filter := domain.AppFilter{
		CategoryID:     categoryID,
		TargetPengguna: targetPengguna,
		StatusLayanan:  statusLayanan,
		RoleIDs:        roleIDs,
		AktifOnly:      true,
		IsPublicOnly:   claims == nil, // If guest, only public apps
		UserID:         userID,
		Offset:         offset,
		Limit:          perPage,
	}

	apps, totalItems, err := u.appRepo.List(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch app catalog: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	meta := &response.Meta{
		Page:       page,
		PerPage:    perPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return apps, meta, nil
}

func (u *CatalogUsecase) GetAppBySlug(ctx context.Context, slug string, userID *int) (*domain.App, error) {
	return u.appRepo.GetBySlug(ctx, slug, userID)
}

func (u *CatalogUsecase) SearchApps(
	ctx context.Context,
	query string,
	categorySlug string,
	claims *jwt.SessionClaims,
	limit int,
) ([]domain.App, error) {
	var categoryID *int
	if categorySlug != "" {
		cat, err := u.categoryRepo.GetBySlug(ctx, categorySlug)
		if err == nil && cat != nil {
			categoryID = &cat.ID
		}
	}

	var userID *int
	var roleIDs []int
	if claims != nil {
		userID = &claims.UserID
		user, err := u.userRepo.GetByID(ctx, claims.UserID)
		if err == nil && user != nil {
			for _, r := range user.Roles {
				roleIDs = append(roleIDs, r.ID)
			}
		}
	}

	isPublicOnly := claims == nil
	return u.appRepo.Search(ctx, query, categoryID, roleIDs, isPublicOnly, userID, limit)
}
