package postgres

import (
	"context"
	"fmt"
	"time"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusCheckRepo struct {
	pool *pgxpool.Pool
}

func NewStatusCheckRepo(pool *pgxpool.Pool) *StatusCheckRepo {
	return &StatusCheckRepo{pool: pool}
}

func (r *StatusCheckRepo) Record(ctx context.Context, check *domain.StatusCheck) error {
	query := `
	INSERT INTO status_checks (app_id, status, http_status_code, latency_ms, error_message, checked_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(
		ctx, query,
		check.AppID, check.Status, check.HTTPStatusCode, check.LatencyMS, check.ErrorMessage, check.CheckedAt,
	)
	return err
}

func (r *StatusCheckRepo) GetRecentByApp(ctx context.Context, appID int, limit int) ([]domain.StatusCheck, error) {
	if limit <= 0 {
		limit = 30
	}

	query := `
	SELECT id, app_id, status, http_status_code, latency_ms, error_message, checked_at
	FROM status_checks
	WHERE app_id = $1
	ORDER BY checked_at DESC
	LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, appID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get status checks: %w", err)
	}
	defer rows.Close()

	var checks []domain.StatusCheck
	for rows.Next() {
		var c domain.StatusCheck
		err := rows.Scan(
			&c.ID, &c.AppID, &c.Status, &c.HTTPStatusCode, &c.LatencyMS, &c.ErrorMessage, &c.CheckedAt,
		)
		if err != nil {
			return nil, err
		}
		checks = append(checks, c)
	}

	return checks, nil
}

func (r *StatusCheckRepo) GetUptimeSummary(ctx context.Context, days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}

	sinceDate := time.Now().AddDate(0, 0, -days)

	query := `
	SELECT 
		a.id, a.nama, a.slug, a.status_layanan,
		COUNT(sc.id) as total_checks,
		COUNT(sc.id) FILTER (WHERE sc.status IN ('online', 'degraded')) as healthy_checks,
		AVG(sc.latency_ms) FILTER (WHERE sc.status IN ('online', 'degraded')) as avg_latency_ms
	FROM apps a
	LEFT JOIN status_checks sc ON a.id = sc.app_id AND sc.checked_at >= $1
	WHERE a.aktif = true
	GROUP BY a.id, a.nama, a.slug, a.status_layanan
	ORDER BY a.urutan ASC, a.nama ASC
	`

	rows, err := r.pool.Query(ctx, query, sinceDate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate uptime summary: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var nama, slug, statusLayanan string
		var totalChecks, healthyChecks int64
		var avgLatency *float64

		err := rows.Scan(&id, &nama, &slug, &statusLayanan, &totalChecks, &healthyChecks, &avgLatency)
		if err != nil {
			return nil, err
		}

		uptimePct := 100.0
		if totalChecks > 0 {
			uptimePct = (float64(healthyChecks) / float64(totalChecks)) * 100.0
		}

		avgLatencyMS := 0
		if avgLatency != nil {
			avgLatencyMS = int(*avgLatency)
		}

		results = append(results, map[string]interface{}{
			"app_id":          id,
			"nama":            nama,
			"slug":            slug,
			"status_layanan":  statusLayanan,
			"uptime_pct":      fmt.Sprintf("%.2f", uptimePct),
			"avg_latency_ms":  avgLatencyMS,
			"total_checks":    totalChecks,
		})
	}

	return results, nil
}
