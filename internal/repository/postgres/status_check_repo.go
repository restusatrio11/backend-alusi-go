package postgres

import (
	"context"
	"fmt"
	"math"
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

func (r *StatusCheckRepo) GetUptimeSummary(ctx context.Context, days int) (*domain.ServiceUptimeSummary, error) {
	summary := &domain.ServiceUptimeSummary{
		OverallUptimePercentage: 100.0,
		Services:                make([]domain.AppUptimeSummary, 0),
	}

	if r == nil || r.pool == nil {
		return summary, nil
	}

	if days <= 0 {
		days = 30
	}

	sinceDate := time.Now().AddDate(0, 0, -days)

	query := `
	SELECT 
		a.id, a.nama, a.slug, a.ikon_url, a.status_layanan,
		COUNT(sc.id) as total_checks,
		COUNT(sc.id) FILTER (WHERE sc.status IN ('online', 'degraded')) as healthy_checks,
		COUNT(sc.id) FILTER (WHERE sc.status IN ('down', 'offline')) as failed_checks,
		AVG(sc.latency_ms) FILTER (WHERE sc.status IN ('online', 'degraded')) as avg_latency_ms,
		MAX(sc.checked_at) as last_check_at
	FROM apps a
	LEFT JOIN status_checks sc ON a.id = sc.app_id AND sc.checked_at >= $1
	WHERE a.aktif = true
	GROUP BY a.id, a.nama, a.slug, a.ikon_url, a.status_layanan
	ORDER BY a.urutan ASC, a.nama ASC
	`

	rows, err := r.pool.Query(ctx, query, sinceDate)
	if err != nil {
		return summary, fmt.Errorf("failed to calculate uptime summary: %w", err)
	}
	defer rows.Close()

	var totalLatencySum int
	var latencyCount int
	var uptimeSum float64

	for rows.Next() {
		var id int
		var nama, slug string
		var ikonURL *string
		var statusLayanan string
		var totalChecks, healthyChecks, failedChecks int64
		var avgLatency *float64
		var lastCheckAt *time.Time

		err := rows.Scan(&id, &nama, &slug, &ikonURL, &statusLayanan, &totalChecks, &healthyChecks, &failedChecks, &avgLatency, &lastCheckAt)
		if err != nil {
			return summary, err
		}

		// Normalize current status
		currentStatus := statusLayanan
		switch statusLayanan {
		case "maintenance", "pemeliharaan":
			currentStatus = "pemeliharaan"
			summary.MaintenanceApps++
		case "degraded", "kendala":
			currentStatus = "kendala"
			summary.DegradedApps++
		case "down", "offline":
			currentStatus = "offline"
			summary.OfflineApps++
		default:
			currentStatus = "online"
			summary.OnlineApps++
		}

		uptimePct := 100.0
		if totalChecks > 0 {
			uptimePct = math.Round((float64(healthyChecks)/float64(totalChecks))*1000) / 10
		} else if currentStatus != "online" {
			uptimePct = 0.0
		}
		uptimeSum += uptimePct

		avgLatencyMS := 0
		if avgLatency != nil {
			avgLatencyMS = int(*avgLatency)
			totalLatencySum += avgLatencyMS
			latencyCount++
		}

		var lastCheckStr *string
		if lastCheckAt != nil {
			s := lastCheckAt.Format(time.RFC3339)
			lastCheckStr = &s
		}

		summary.Services = append(summary.Services, domain.AppUptimeSummary{
			AppID:            id,
			AppNama:          nama,
			AppSlug:          slug,
			AppIkonURL:       ikonURL,
			CurrentStatus:    currentStatus,
			UptimePercentage: uptimePct,
			AvgLatencyMS:     avgLatencyMS,
			LastCheckAt:      lastCheckStr,
			TotalChecks:      totalChecks,
			SuccessChecks:    healthyChecks,
			FailedChecks:     failedChecks,
		})
	}

	summary.TotalApps = len(summary.Services)
	if summary.TotalApps > 0 {
		summary.OverallUptimePercentage = math.Round((uptimeSum/float64(summary.TotalApps))*100) / 100
	}
	if latencyCount > 0 {
		summary.AvgResponseTimeMS = totalLatencySum / latencyCount
	}

	return summary, nil
}
