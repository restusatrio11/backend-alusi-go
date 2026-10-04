package usecase

import (
	"context"
	"fmt"
	"strconv"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
	"backend-alusi-go/pkg/worker"
)

type MonitoringUsecase struct {
	statusCheckRepo   *postgres.StatusCheckRepo
	appRepo           *postgres.AppRepo
	healthProbeWorker *worker.HealthProbeWorker
}

func NewMonitoringUsecase(
	statusCheckRepo *postgres.StatusCheckRepo,
	appRepo *postgres.AppRepo,
	healthProbeWorker *worker.HealthProbeWorker,
) *MonitoringUsecase {
	return &MonitoringUsecase{
		statusCheckRepo:   statusCheckRepo,
		appRepo:           appRepo,
		healthProbeWorker: healthProbeWorker,
	}
}

// GetAppStatusHistory returns historical health check probes for an application
func (u *MonitoringUsecase) GetAppStatusHistory(ctx context.Context, appID int, limit int) ([]domain.StatusCheck, error) {
	if u.statusCheckRepo == nil {
		return []domain.StatusCheck{}, nil
	}
	return u.statusCheckRepo.GetRecentByApp(ctx, appID, limit)
}

// GetAppStatusHistoryByIdentifier resolves slug or ID then returns history
func (u *MonitoringUsecase) GetAppStatusHistoryByIdentifier(ctx context.Context, identifier string, limit int) ([]domain.StatusCheck, error) {
	if u.statusCheckRepo == nil {
		return []domain.StatusCheck{}, nil
	}

	appID := 0
	if id, err := strconv.Atoi(identifier); err == nil {
		appID = id
	} else if u.appRepo != nil {
		app, err := u.appRepo.GetBySlug(ctx, identifier, nil)
		if err != nil || app == nil {
			return nil, fmt.Errorf("aplikasi '%s' tidak ditemukan", identifier)
		}
		appID = app.ID
	}

	return u.statusCheckRepo.GetRecentByApp(ctx, appID, limit)
}

// GetServiceUptimeSummary returns 30-day uptime summary for all applications
func (u *MonitoringUsecase) GetServiceUptimeSummary(ctx context.Context, days int) (*domain.ServiceUptimeSummary, error) {
	if u.statusCheckRepo == nil {
		return &domain.ServiceUptimeSummary{
			OverallUptimePercentage: 100.0,
			Services:                make([]domain.AppUptimeSummary, 0),
		}, nil
	}
	return u.statusCheckRepo.GetUptimeSummary(ctx, days)
}

// ManualProbeApp triggers an on-demand probe for an application
func (u *MonitoringUsecase) ManualProbeApp(ctx context.Context, appID int) (*domain.StatusCheck, error) {
	app, err := u.appRepo.GetByID(ctx, appID, nil)
	if err != nil || app == nil {
		return nil, fmt.Errorf("aplikasi tidak ditemukan")
	}

	if u.healthProbeWorker != nil {
		return u.healthProbeWorker.ProbeSingleApp(ctx, app)
	}

	return nil, fmt.Errorf("probe worker not initialized")
}
