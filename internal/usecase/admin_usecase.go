package usecase

import (
	"context"
	"fmt"
	"strings"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
)

type AdminUsecase struct {
	appRepo      *postgres.AppRepo
	categoryRepo *postgres.CategoryRepo
	guideRepo    *postgres.GuideRepo
}

func NewAdminUsecase(
	appRepo *postgres.AppRepo,
	categoryRepo *postgres.CategoryRepo,
	guideRepo *postgres.GuideRepo,
) *AdminUsecase {
	return &AdminUsecase{
		appRepo:      appRepo,
		categoryRepo: categoryRepo,
		guideRepo:    guideRepo,
	}
}

type CreateAppInput struct {
	CategoryID     int     `json:"category_id" binding:"required"`
	Nama           string  `json:"nama" binding:"required"`
	Slug           string  `json:"slug"`
	URL            string  `json:"url" binding:"required"`
	Deskripsi      *string `json:"deskripsi"`
	IkonURL        *string `json:"ikon_url"`
	TargetPengguna string  `json:"target_pengguna"`
	Pemilik        string  `json:"pemilik"`
	KontakAdmin    *string `json:"kontak_admin"`
	Urutan         int     `json:"urutan"`
	StatusLayanan  string  `json:"status_layanan"`
	IsPublic       bool    `json:"is_public"`
	RoleIDs        []int   `json:"role_ids"`
}

type UpdateAppInput struct {
	CategoryID     int     `json:"category_id" binding:"required"`
	Nama           string  `json:"nama" binding:"required"`
	Slug           string  `json:"slug"`
	URL            string  `json:"url" binding:"required"`
	Deskripsi      *string `json:"deskripsi"`
	IkonURL        *string `json:"ikon_url"`
	TargetPengguna string  `json:"target_pengguna"`
	Pemilik        string  `json:"pemilik"`
	KontakAdmin    *string `json:"kontak_admin"`
	Urutan         int     `json:"urutan"`
	StatusLayanan  string  `json:"status_layanan"`
	IsPublic       bool    `json:"is_public"`
	Aktif          bool    `json:"aktif"`
	RoleIDs        []int   `json:"role_ids"`
}

func (u *AdminUsecase) CreateApp(ctx context.Context, input CreateAppInput) (*domain.App, error) {
	slug := input.Slug
	if slug == "" {
		slug = generateSlug(input.Nama)
	}

	target := input.TargetPengguna
	if target == "" {
		target = "Semua"
	}

	pemilik := input.Pemilik
	if pemilik == "" {
		pemilik = "BPS RI"
	}

	status := input.StatusLayanan
	if status == "" {
		status = "online"
	}

	app := &domain.App{
		CategoryID:     input.CategoryID,
		Nama:           input.Nama,
		Slug:           slug,
		URL:            input.URL,
		Deskripsi:      input.Deskripsi,
		IkonURL:        input.IkonURL,
		TargetPengguna: target,
		Pemilik:        pemilik,
		KontakAdmin:    input.KontakAdmin,
		Urutan:         input.Urutan,
		StatusLayanan:  status,
		IsPublic:       input.IsPublic,
		Aktif:          true,
	}

	if err := u.appRepo.Create(ctx, app, input.RoleIDs); err != nil {
		return nil, fmt.Errorf("gagal membuat aplikasi: %w", err)
	}

	return app, nil
}

func (u *AdminUsecase) UpdateApp(ctx context.Context, id int, input UpdateAppInput) (*domain.App, error) {
	app, err := u.appRepo.GetByID(ctx, id, nil)
	if err != nil || app == nil {
		return nil, fmt.Errorf("aplikasi tidak ditemukan")
	}

	slug := input.Slug
	if slug == "" {
		slug = generateSlug(input.Nama)
	}

	app.CategoryID = input.CategoryID
	app.Nama = input.Nama
	app.Slug = slug
	app.URL = input.URL
	app.Deskripsi = input.Deskripsi
	app.IkonURL = input.IkonURL
	app.TargetPengguna = input.TargetPengguna
	app.Pemilik = input.Pemilik
	app.KontakAdmin = input.KontakAdmin
	app.Urutan = input.Urutan
	app.StatusLayanan = input.StatusLayanan
	app.IsPublic = input.IsPublic
	app.Aktif = input.Aktif

	if err := u.appRepo.Update(ctx, app, input.RoleIDs); err != nil {
		return nil, fmt.Errorf("gagal memperbarui aplikasi: %w", err)
	}

	return app, nil
}

func (u *AdminUsecase) DeleteApp(ctx context.Context, id int) error {
	return u.appRepo.Delete(ctx, id)
}

func (u *AdminUsecase) ReorderApps(ctx context.Context, appIDs []int) error {
	return u.appRepo.Reorder(ctx, appIDs)
}

type CategoryInput struct {
	Nama      string  `json:"nama" binding:"required"`
	Slug      string  `json:"slug"`
	Deskripsi *string `json:"deskripsi"`
	Urutan    int     `json:"urutan"`
	Ikon      *string `json:"ikon"`
}

func (u *AdminUsecase) CreateCategory(ctx context.Context, input CategoryInput) (*domain.Category, error) {
	slug := input.Slug
	if slug == "" {
		slug = generateSlug(input.Nama)
	}

	cat := &domain.Category{
		Nama:      input.Nama,
		Slug:      slug,
		Deskripsi: input.Deskripsi,
		Urutan:    input.Urutan,
		Ikon:      input.Ikon,
	}

	if err := u.categoryRepo.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("gagal membuat kategori: %w", err)
	}

	return cat, nil
}

func (u *AdminUsecase) UpdateCategory(ctx context.Context, id int, input CategoryInput) (*domain.Category, error) {
	cat, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil || cat == nil {
		return nil, fmt.Errorf("kategori tidak ditemukan")
	}

	slug := input.Slug
	if slug == "" {
		slug = generateSlug(input.Nama)
	}

	cat.Nama = input.Nama
	cat.Slug = slug
	cat.Deskripsi = input.Deskripsi
	cat.Urutan = input.Urutan
	cat.Ikon = input.Ikon

	if err := u.categoryRepo.Update(ctx, cat); err != nil {
		return nil, fmt.Errorf("gagal memperbarui kategori: %w", err)
	}

	return cat, nil
}

func (u *AdminUsecase) DeleteCategory(ctx context.Context, id int) error {
	return u.categoryRepo.Delete(ctx, id)
}

type GuideInput struct {
	Judul  string `json:"judul" binding:"required"`
	Konten string `json:"konten" binding:"required"`
	Urutan int    `json:"urutan"`
}

func (u *AdminUsecase) CreateGuide(ctx context.Context, appID int, input GuideInput) (*domain.Guide, error) {
	guide := &domain.Guide{
		AppID:  appID,
		Judul:  input.Judul,
		Konten: input.Konten,
		Urutan: input.Urutan,
	}

	if err := u.guideRepo.Create(ctx, guide); err != nil {
		return nil, fmt.Errorf("gagal menambahkan panduan: %w", err)
	}

	return guide, nil
}

func (u *AdminUsecase) UpdateGuide(ctx context.Context, id int, input GuideInput) (*domain.Guide, error) {
	guide, err := u.guideRepo.GetByID(ctx, id)
	if err != nil || guide == nil {
		return nil, fmt.Errorf("panduan tidak ditemukan")
	}

	guide.Judul = input.Judul
	guide.Konten = input.Konten
	guide.Urutan = input.Urutan

	if err := u.guideRepo.Update(ctx, guide); err != nil {
		return nil, fmt.Errorf("gagal memperbarui panduan: %w", err)
	}

	return guide, nil
}

func (u *AdminUsecase) DeleteGuide(ctx context.Context, id int) error {
	return u.guideRepo.Delete(ctx, id)
}

func generateSlug(text string) string {
	slug := strings.ToLower(text)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, "&", "dan")
	return slug
}
