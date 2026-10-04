package postgres

import (
	"context"
	"fmt"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GuideRepo struct {
	pool *pgxpool.Pool
}

func NewGuideRepo(pool *pgxpool.Pool) *GuideRepo {
	return &GuideRepo{pool: pool}
}

func (r *GuideRepo) ListByAppID(ctx context.Context, appID int) ([]domain.Guide, error) {
	query := `
	SELECT id, app_id, judul, konten, urutan, created_at, updated_at
	FROM guides
	WHERE app_id = $1
	ORDER BY urutan ASC, id ASC
	`
	rows, err := r.pool.Query(ctx, query, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to list guides: %w", err)
	}
	defer rows.Close()

	var guides []domain.Guide
	for rows.Next() {
		var g domain.Guide
		if err := rows.Scan(&g.ID, &g.AppID, &g.Judul, &g.Konten, &g.Urutan, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		guides = append(guides, g)
	}

	return guides, nil
}

func (r *GuideRepo) GetByID(ctx context.Context, id int) (*domain.Guide, error) {
	query := `
	SELECT id, app_id, judul, konten, urutan, created_at, updated_at
	FROM guides
	WHERE id = $1
	`
	var g domain.Guide
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&g.ID, &g.AppID, &g.Judul, &g.Konten, &g.Urutan, &g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *GuideRepo) Create(ctx context.Context, guide *domain.Guide) error {
	query := `
	INSERT INTO guides (app_id, judul, konten, urutan)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(
		ctx, query,
		guide.AppID, guide.Judul, guide.Konten, guide.Urutan,
	).Scan(&guide.ID, &guide.CreatedAt, &guide.UpdatedAt)
}

func (r *GuideRepo) Update(ctx context.Context, guide *domain.Guide) error {
	query := `
	UPDATE guides
	SET judul = $1, konten = $2, urutan = $3, updated_at = NOW()
	WHERE id = $4
	RETURNING updated_at
	`
	return r.pool.QueryRow(
		ctx, query,
		guide.Judul, guide.Konten, guide.Urutan, guide.ID,
	).Scan(&guide.UpdatedAt)
}

func (r *GuideRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM guides WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
