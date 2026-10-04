package postgres

import (
	"context"
	"fmt"
	"time"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnnouncementRepo struct {
	pool *pgxpool.Pool
}

func NewAnnouncementRepo(pool *pgxpool.Pool) *AnnouncementRepo {
	return &AnnouncementRepo{pool: pool}
}

// ListActive returns all currently active announcements where NOW() between starts_at and COALESCE(ends_at, infinity)
func (r *AnnouncementRepo) ListActive(ctx context.Context, appID *int) ([]domain.Announcement, error) {
	whereClause := "WHERE a.is_active = true AND a.starts_at <= NOW() AND (a.ends_at IS NULL OR a.ends_at >= NOW())"
	var args []interface{}

	if appID != nil {
		whereClause += " AND (a.app_id = $1 OR a.app_id IS NULL)"
		args = append(args, *appID)
	}

	query := fmt.Sprintf(`
	SELECT 
		a.id, a.judul, a.konten, a.tipe, a.app_id, a.is_active, a.starts_at, a.ends_at, a.created_at, a.updated_at,
		ap.id, ap.nama, ap.slug
	FROM announcements a
	LEFT JOIN apps ap ON a.app_id = ap.id
	%s
	ORDER BY 
		CASE WHEN a.tipe = 'urgent' THEN 1 WHEN a.tipe = 'maintenance' THEN 2 WHEN a.tipe = 'warning' THEN 3 ELSE 4 END,
		a.created_at DESC
	`, whereClause)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list active announcements: %w", err)
	}
	defer rows.Close()

	var list []domain.Announcement
	for rows.Next() {
		ann, err := r.scanAnnouncementRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *ann)
	}

	return list, nil
}

// ListAll returns paginated announcements for admin
func (r *AnnouncementRepo) ListAll(ctx context.Context, limit, offset int) ([]domain.Announcement, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM announcements`
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count announcements: %w", err)
	}

	query := `
	SELECT 
		a.id, a.judul, a.konten, a.tipe, a.app_id, a.is_active, a.starts_at, a.ends_at, a.created_at, a.updated_at,
		ap.id, ap.nama, ap.slug
	FROM announcements a
	LEFT JOIN apps ap ON a.app_id = ap.id
	ORDER BY a.created_at DESC
	LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list announcements: %w", err)
	}
	defer rows.Close()

	var list []domain.Announcement
	for rows.Next() {
		ann, err := r.scanAnnouncementRow(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, *ann)
	}

	return list, total, nil
}

func (r *AnnouncementRepo) GetByID(ctx context.Context, id int) (*domain.Announcement, error) {
	query := `
	SELECT 
		a.id, a.judul, a.konten, a.tipe, a.app_id, a.is_active, a.starts_at, a.ends_at, a.created_at, a.updated_at,
		ap.id, ap.nama, ap.slug
	FROM announcements a
	LEFT JOIN apps ap ON a.app_id = ap.id
	WHERE a.id = $1
	`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	return r.scanAnnouncementRow(rows)
}

func (r *AnnouncementRepo) Create(ctx context.Context, ann *domain.Announcement) error {
	query := `
	INSERT INTO announcements (judul, konten, tipe, app_id, is_active, starts_at, ends_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at
	`
	if ann.StartsAt.IsZero() {
		ann.StartsAt = time.Now()
	}
	return r.pool.QueryRow(
		ctx, query,
		ann.Judul, ann.Konten, ann.Tipe, ann.AppID, ann.IsActive, ann.StartsAt, ann.EndsAt,
	).Scan(&ann.ID, &ann.CreatedAt, &ann.UpdatedAt)
}

func (r *AnnouncementRepo) Update(ctx context.Context, ann *domain.Announcement) error {
	query := `
	UPDATE announcements
	SET judul = $1, konten = $2, tipe = $3, app_id = $4, is_active = $5, starts_at = $6, ends_at = $7, updated_at = NOW()
	WHERE id = $8
	RETURNING updated_at
	`
	return r.pool.QueryRow(
		ctx, query,
		ann.Judul, ann.Konten, ann.Tipe, ann.AppID, ann.IsActive, ann.StartsAt, ann.EndsAt, ann.ID,
	).Scan(&ann.UpdatedAt)
}

func (r *AnnouncementRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM announcements WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *AnnouncementRepo) scanAnnouncementRow(rows pgx.Rows) (*domain.Announcement, error) {
	var a domain.Announcement
	var appID *int
	var appNama *string
	var appSlug *string

	err := rows.Scan(
		&a.ID, &a.Judul, &a.Konten, &a.Tipe, &a.AppID, &a.IsActive, &a.StartsAt, &a.EndsAt, &a.CreatedAt, &a.UpdatedAt,
		&appID, &appNama, &appSlug,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan announcement: %w", err)
	}

	if appID != nil && appNama != nil && appSlug != nil {
		a.App = &domain.App{
			ID:   *appID,
			Nama: *appNama,
			Slug: *appSlug,
		}
	}

	return &a, nil
}
