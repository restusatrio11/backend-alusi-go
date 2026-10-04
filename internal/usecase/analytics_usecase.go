package usecase

import (
	"context"

	"backend-alusi-go/internal/repository/postgres"
)

type AnalyticsUsecase struct {
	analyticsRepo *postgres.AnalyticsRepo
}

func NewAnalyticsUsecase(analyticsRepo *postgres.AnalyticsRepo) *AnalyticsUsecase {
	return &AnalyticsUsecase{
		analyticsRepo: analyticsRepo,
	}
}

func (u *AnalyticsUsecase) GetDashboardSummary(ctx context.Context) (*postgres.DashboardSummary, error) {
	if u.analyticsRepo == nil {
		return &postgres.DashboardSummary{}, nil
	}
	return u.analyticsRepo.GetDashboardSummary(ctx)
}

func (u *AnalyticsUsecase) GetTopApps(ctx context.Context, days int, limit int) ([]postgres.TopAppMetric, error) {
	if u.analyticsRepo == nil {
		return []postgres.TopAppMetric{}, nil
	}
	return u.analyticsRepo.GetTopApps(ctx, days, limit)
}

func (u *AnalyticsUsecase) GetDailyClicksTrend(ctx context.Context, days int) ([]postgres.DailyClickTrend, error) {
	if u.analyticsRepo == nil {
		return []postgres.DailyClickTrend{}, nil
	}
	return u.analyticsRepo.GetDailyClicksTrend(ctx, days)
}

func (u *AnalyticsUsecase) GetDisruptionSummary(ctx context.Context, days int) ([]postgres.DisruptionSummary, error) {
	if u.analyticsRepo == nil {
		return []postgres.DisruptionSummary{}, nil
	}
	return u.analyticsRepo.GetDisruptionSummary(ctx, days)
}
