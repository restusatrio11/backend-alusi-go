package postgres

import (
	"context"
	"fmt"
	"strings"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedbackRepo struct {
	pool *pgxpool.Pool
}

func NewFeedbackRepo(pool *pgxpool.Pool) *FeedbackRepo {
	return &FeedbackRepo{pool: pool}
}

func (r *FeedbackRepo) Create(ctx context.Context, fb *domain.Feedback) error {
	query := `
	INSERT INTO feedbacks (user_id, app_id, kategori, pesan, status)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	status := fb.Status
	if status == "" {
		status = "pending"
	}
	kategori := fb.Kategori
	if kategori == "" {
		kategori = "kendala"
	}

	return r.pool.QueryRow(
		ctx, query,
		fb.UserID, fb.AppID, kategori, fb.Pesan, status,
	).Scan(&fb.ID, &fb.CreatedAt, &fb.UpdatedAt)
}

func (r *FeedbackRepo) List(ctx context.Context, status string, appID *int, limit, offset int) ([]domain.Feedback, int64, error) {
	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("f.status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}

	if appID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("f.app_id = $%d", argIdx))
		args = append(args, *appID)
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM feedbacks f %s", whereSQL)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count feedbacks: %w", err)
	}

	limitSQL := fmt.Sprintf("LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
	SELECT 
		f.id, f.user_id, f.app_id, f.kategori, f.pesan, f.status, f.catatan_admin, f.created_at, f.updated_at,
		u.id, u.nama, u.email,
		a.id, a.nama, a.slug
	FROM feedbacks f
	LEFT JOIN users u ON f.user_id = u.id
	LEFT JOIN apps a ON f.app_id = a.id
	%s
	ORDER BY f.created_at DESC
	%s
	`, whereSQL, limitSQL)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list feedbacks: %w", err)
	}
	defer rows.Close()

	var list []domain.Feedback
	for rows.Next() {
		fb, err := r.scanFeedbackRow(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, *fb)
	}

	return list, total, nil
}

func (r *FeedbackRepo) GetByID(ctx context.Context, id int) (*domain.Feedback, error) {
	query := `
	SELECT 
		f.id, f.user_id, f.app_id, f.kategori, f.pesan, f.status, f.catatan_admin, f.created_at, f.updated_at,
		u.id, u.nama, u.email,
		a.id, a.nama, a.slug
	FROM feedbacks f
	LEFT JOIN users u ON f.user_id = u.id
	LEFT JOIN apps a ON f.app_id = a.id
	WHERE f.id = $1
	`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	return r.scanFeedbackRow(rows)
}

func (r *FeedbackRepo) UpdateStatus(ctx context.Context, id int, status string, catatanAdmin *string) error {
	query := `
	UPDATE feedbacks
	SET status = $1, catatan_admin = $2, updated_at = NOW()
	WHERE id = $3
	`
	_, err := r.pool.Exec(ctx, query, status, catatanAdmin, id)
	return err
}

func (r *FeedbackRepo) scanFeedbackRow(rows pgx.Rows) (*domain.Feedback, error) {
	var f domain.Feedback
	var uID *int
	var uNama, uEmail *string
	var aID *int
	var aNama, aSlug *string

	err := rows.Scan(
		&f.ID, &f.UserID, &f.AppID, &f.Kategori, &f.Pesan, &f.Status, &f.CatatanAdmin, &f.CreatedAt, &f.UpdatedAt,
		&uID, &uNama, &uEmail,
		&aID, &aNama, &aSlug,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan feedback row: %w", err)
	}

	if uID != nil && uNama != nil && uEmail != nil {
		f.User = &domain.User{
			ID:    *uID,
			Nama:  *uNama,
			Email: *uEmail,
		}
	}

	if aID != nil && aNama != nil && aSlug != nil {
		f.App = &domain.App{
			ID:   *aID,
			Nama: *aNama,
			Slug: *aSlug,
		}
	}

	return &f, nil
}
