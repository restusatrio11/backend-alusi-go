package postgres

import (
	"context"
	"fmt"
	"time"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClickLogRepo struct {
	pool *pgxpool.Pool
}

func NewClickLogRepo(pool *pgxpool.Pool) *ClickLogRepo {
	return &ClickLogRepo{pool: pool}
}

func (r *ClickLogRepo) Record(ctx context.Context, log *domain.ClickLog) error {
	query := `
	INSERT INTO click_logs (user_id, app_id, satker_id, ip_address, user_agent, clicked_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(
		ctx, query,
		log.UserID, log.AppID, log.SatkerID, log.IPAddress, log.UserAgent, log.ClickedAt,
	)
	return err
}

func (r *ClickLogRepo) GetRecentByUser(ctx context.Context, userID int, limit int) ([]domain.App, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
	WITH recent_apps AS (
		SELECT DISTINCT ON (cl.app_id) cl.app_id, cl.clicked_at
		FROM click_logs cl
		WHERE cl.user_id = $1
		ORDER BY cl.app_id, cl.clicked_at DESC
	)
	SELECT 
		a.id, a.category_id, a.nama, a.slug, a.url, a.deskripsi, a.ikon_url,
		a.target_pengguna, a.pemilik, a.kontak_admin, a.urutan, a.status_layanan,
		a.is_public, a.aktif, a.created_at, a.updated_at,
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at,
		(f.user_id IS NOT NULL) as is_favorite
	FROM recent_apps ra
	INNER JOIN apps a ON ra.app_id = a.id
	INNER JOIN categories c ON a.category_id = c.id
	LEFT JOIN favorites f ON a.id = f.app_id AND f.user_id = $1
	WHERE a.aktif = true
	ORDER BY ra.clicked_at DESC
	LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user recent apps: %w", err)
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

func (r *ClickLogRepo) GetTopApps(ctx context.Context, startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
	SELECT 
		a.id, a.nama, a.slug, a.url, a.ikon_url, c.nama as category_nama,
		COUNT(cl.id) as click_count
	FROM click_logs cl
	INNER JOIN apps a ON cl.app_id = a.id
	INNER JOIN categories c ON a.category_id = c.id
	WHERE cl.clicked_at >= $1 AND cl.clicked_at <= $2 AND a.aktif = true
	GROUP BY a.id, a.nama, a.slug, a.url, a.ikon_url, c.nama
	ORDER BY click_count DESC
	LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top apps: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var nama, slug, appURL, categoryNama string
		var ikonURL *string
		var clickCount int64

		err := rows.Scan(&id, &nama, &slug, &appURL, &ikonURL, &categoryNama, &clickCount)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"id":            id,
			"nama":          nama,
			"slug":          slug,
			"url":           appURL,
			"ikon_url":      ikonURL,
			"category_nama": categoryNama,
			"click_count":   clickCount,
		})
	}

	return results, nil
}

func (r *ClickLogRepo) GetSummary(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	query := `
	SELECT 
		COUNT(id) as total_clicks,
		COUNT(DISTINCT user_id) as total_active_users,
		COUNT(DISTINCT app_id) as total_active_apps
	FROM click_logs
	WHERE clicked_at >= $1 AND clicked_at <= $2
	`
	var totalClicks, totalActiveUsers, totalActiveApps int64
	err := r.pool.QueryRow(ctx, query, startDate, endDate).Scan(&totalClicks, &totalActiveUsers, &totalActiveApps)
	if err != nil {
		return nil, fmt.Errorf("failed to get click summary: %w", err)
	}

	return map[string]interface{}{
		"total_clicks":       totalClicks,
		"total_active_users": totalActiveUsers,
		"total_active_apps":  totalActiveApps,
	}, nil
}
