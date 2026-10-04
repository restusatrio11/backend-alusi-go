package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepo struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepo(pool *pgxpool.Pool) *AnalyticsRepo {
	return &AnalyticsRepo{pool: pool}
}

type DashboardSummary struct {
	TotalApps             int64 `json:"total_apps"`
	TotalCategories       int64 `json:"total_categories"`
	TotalUsers            int64 `json:"total_users"`
	TotalClicksToday      int64 `json:"total_clicks_today"`
	TotalClicksThisMonth  int64 `json:"total_clicks_this_month"`
	DAU                   int64 `json:"dau"` // Daily Active Users
	MAU                   int64 `json:"mau"` // Monthly Active Users
	PendingFeedbacksCount int64 `json:"pending_feedbacks_count"`
}

func (r *AnalyticsRepo) GetDashboardSummary(ctx context.Context) (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	// 1. Total Apps
	_ = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM apps WHERE aktif = true").Scan(&summary.TotalApps)

	// 2. Total Categories
	_ = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM categories").Scan(&summary.TotalCategories)

	// 3. Total Users
	_ = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE status = 'active'").Scan(&summary.TotalUsers)

	// 4. Total Clicks Today & DAU
	todayQuery := `
	SELECT 
		COUNT(*), 
		COUNT(DISTINCT user_id) 
	FROM click_logs 
	WHERE clicked_at >= CURRENT_DATE
	`
	_ = r.pool.QueryRow(ctx, todayQuery).Scan(&summary.TotalClicksToday, &summary.DAU)

	// 5. Total Clicks This Month & MAU
	monthQuery := `
	SELECT 
		COUNT(*), 
		COUNT(DISTINCT user_id) 
	FROM click_logs 
	WHERE clicked_at >= date_trunc('month', CURRENT_DATE)
	`
	_ = r.pool.QueryRow(ctx, monthQuery).Scan(&summary.TotalClicksThisMonth, &summary.MAU)

	// 6. Pending Feedbacks
	_ = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM feedbacks WHERE status = 'pending'").Scan(&summary.PendingFeedbacksCount)

	return summary, nil
}

type TopAppMetric struct {
	AppID       int    `json:"app_id"`
	Nama        string `json:"nama"`
	Slug        string `json:"slug"`
	Category    string `json:"category"`
	TotalClicks int64  `json:"total_clicks"`
	UniqueUsers int64  `json:"unique_users"`
}

func (r *AnalyticsRepo) GetTopApps(ctx context.Context, days int, limit int) ([]TopAppMetric, error) {
	if days <= 0 {
		days = 30
	}
	if limit <= 0 {
		limit = 10
	}

	query := `
	SELECT 
		a.id, a.nama, a.slug, c.nama as category_nama,
		COUNT(cl.id) as total_clicks,
		COUNT(DISTINCT cl.user_id) as unique_users
	FROM apps a
	INNER JOIN categories c ON a.category_id = c.id
	LEFT JOIN click_logs cl ON a.id = cl.app_id AND cl.clicked_at >= NOW() - ($1 || ' days')::INTERVAL
	WHERE a.aktif = true
	GROUP BY a.id, a.nama, a.slug, c.nama
	ORDER BY total_clicks DESC, a.nama ASC
	LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, fmt.Sprintf("%d", days), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top apps: %w", err)
	}
	defer rows.Close()

	var list []TopAppMetric
	for rows.Next() {
		var m TopAppMetric
		if err := rows.Scan(&m.AppID, &m.Nama, &m.Slug, &m.Category, &m.TotalClicks, &m.UniqueUsers); err != nil {
			return nil, err
		}
		list = append(list, m)
	}

	return list, nil
}

type DailyClickTrend struct {
	Date        string `json:"date"`
	TotalClicks int64  `json:"total_clicks"`
	UniqueUsers int64  `json:"unique_users"`
}

func (r *AnalyticsRepo) GetDailyClicksTrend(ctx context.Context, days int) ([]DailyClickTrend, error) {
	if days <= 0 {
		days = 30
	}

	query := `
	SELECT 
		TO_CHAR(DATE(clicked_at), 'YYYY-MM-DD') as click_date,
		COUNT(*) as total_clicks,
		COUNT(DISTINCT user_id) as unique_users
	FROM click_logs
	WHERE clicked_at >= NOW() - ($1 || ' days')::INTERVAL
	GROUP BY DATE(clicked_at)
	ORDER BY click_date ASC
	`

	rows, err := r.pool.Query(ctx, query, fmt.Sprintf("%d", days))
	if err != nil {
		return nil, fmt.Errorf("failed to get click trends: %w", err)
	}
	defer rows.Close()

	var list []DailyClickTrend
	for rows.Next() {
		var t DailyClickTrend
		if err := rows.Scan(&t.Date, &t.TotalClicks, &t.UniqueUsers); err != nil {
			return nil, err
		}
		list = append(list, t)
	}

	return list, nil
}

type DisruptionSummary struct {
	AppID            int        `json:"app_id"`
	AppNama          string     `json:"app_nama"`
	AppSlug          string     `json:"app_slug"`
	TotalChecks      int64      `json:"total_checks"`
	DownChecks       int64      `json:"down_checks"`
	DegradedChecks   int64      `json:"degraded_checks"`
	UptimePercentage float64    `json:"uptime_percentage"`
	LastDownAt       *time.Time `json:"last_down_at,omitempty"`
}

func (r *AnalyticsRepo) GetDisruptionSummary(ctx context.Context, days int) ([]DisruptionSummary, error) {
	if days <= 0 {
		days = 30
	}

	query := `
	SELECT 
		a.id, a.nama, a.slug,
		COUNT(sc.id) as total_checks,
		COUNT(CASE WHEN sc.status = 'down' THEN 1 END) as down_checks,
		COUNT(CASE WHEN sc.status = 'degraded' THEN 1 END) as degraded_checks,
		CASE 
			WHEN COUNT(sc.id) = 0 THEN 100.0
			ELSE ROUND((COUNT(CASE WHEN sc.status = 'online' THEN 1 END)::NUMERIC / COUNT(sc.id)::NUMERIC) * 100, 2)
		END as uptime_percentage,
		MAX(CASE WHEN sc.status = 'down' THEN sc.checked_at END) as last_down_at
	FROM apps a
	LEFT JOIN status_checks sc ON a.id = sc.app_id AND sc.checked_at >= NOW() - ($1 || ' days')::INTERVAL
	WHERE a.aktif = true
	GROUP BY a.id, a.nama, a.slug
	ORDER BY down_checks DESC, uptime_percentage ASC
	`

	rows, err := r.pool.Query(ctx, query, fmt.Sprintf("%d", days))
	if err != nil {
		return nil, fmt.Errorf("failed to get disruption summary: %w", err)
	}
	defer rows.Close()

	var list []DisruptionSummary
	for rows.Next() {
		var d DisruptionSummary
		if err := rows.Scan(
			&d.AppID, &d.AppNama, &d.AppSlug,
			&d.TotalChecks, &d.DownChecks, &d.DegradedChecks,
			&d.UptimePercentage, &d.LastDownAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, nil
}
