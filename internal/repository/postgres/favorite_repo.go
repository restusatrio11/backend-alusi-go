package postgres

import (
	"context"
	"fmt"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FavoriteRepo struct {
	pool *pgxpool.Pool
}

func NewFavoriteRepo(pool *pgxpool.Pool) *FavoriteRepo {
	return &FavoriteRepo{pool: pool}
}

func (r *FavoriteRepo) Add(ctx context.Context, userID, appID int) error {
	query := `INSERT INTO favorites (user_id, app_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, userID, appID)
	return err
}

func (r *FavoriteRepo) Remove(ctx context.Context, userID, appID int) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND app_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, appID)
	return err
}

func (r *FavoriteRepo) IsFavorite(ctx context.Context, userID, appID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND app_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, appID).Scan(&exists)
	return exists, err
}

func (r *FavoriteRepo) ListByUser(ctx context.Context, userID int) ([]domain.App, error) {
	query := `
	SELECT 
		a.id, a.category_id, a.nama, a.slug, a.url, a.deskripsi, a.ikon_url,
		a.target_pengguna, a.pemilik, a.kontak_admin, a.urutan, a.status_layanan,
		a.is_public, a.aktif, a.created_at, a.updated_at,
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at,
		true as is_favorite
	FROM favorites f
	INNER JOIN apps a ON f.app_id = a.id
	INNER JOIN categories c ON a.category_id = c.id
	WHERE f.user_id = $1 AND a.aktif = true
	ORDER BY f.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user favorites: %w", err)
	}
	defer rows.Close()

	var apps []domain.App
	for rows.Next() {
		var a domain.App
		var c domain.Category
		err := rows.Scan(
			&a.ID, &a.CategoryID, &a.Nama, &a.Slug, &a.URL, &a.Deskripsi, &a.IkonURL,
			&a.TargetPengguna, &a.Pemilik, &a.KontakAdmin, &a.Urutan, &a.StatusLayanan,
			&a.IsPublic, &a.Aktif, &a.CreatedAt, &a.UpdatedAt,
			&c.ID, &c.Nama, &c.Slug, &c.Deskripsi, &c.Urutan, &c.Ikon, &c.CreatedAt, &c.UpdatedAt,
			&a.IsFavorite,
		)
		if err != nil {
			return nil, err
		}
		a.Category = &c
		apps = append(apps, a)
	}

	return apps, nil
}
