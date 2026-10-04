package usecase

import (
	"context"
	"fmt"
	"time"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/worker"
)

type InteractionUsecase struct {
	favRepo      *postgres.FavoriteRepo
	clickLogRepo *postgres.ClickLogRepo
	appRepo      *postgres.AppRepo
	clickWorker  *worker.ClickWorker
}

func NewInteractionUsecase(
	favRepo *postgres.FavoriteRepo,
	clickLogRepo *postgres.ClickLogRepo,
	appRepo *postgres.AppRepo,
	clickWorker *worker.ClickWorker,
) *InteractionUsecase {
	return &InteractionUsecase{
		favRepo:      favRepo,
		clickLogRepo: clickLogRepo,
		appRepo:      appRepo,
		clickWorker:  clickWorker,
	}
}

// RecordClick enqueues a click log asynchronously and returns the application URL
func (u *InteractionUsecase) RecordClick(
	ctx context.Context,
	appID int,
	claims *jwt.SessionClaims,
	ipAddress, userAgent string,
) (*domain.App, error) {
	app, err := u.appRepo.GetByID(ctx, appID, nil)
	if err != nil || app == nil {
		return nil, fmt.Errorf("aplikasi tidak ditemukan")
	}

	var userID *int
	var satkerID *int
	if claims != nil {
		userID = &claims.UserID
	}

	logItem := &domain.ClickLog{
		UserID:    userID,
		AppID:     app.ID,
		SatkerID:  satkerID,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		ClickedAt: time.Now(),
	}

	u.clickWorker.Enqueue(logItem)
	return app, nil
}

// ToggleFavorite adds or removes an app from user's favorites
func (u *InteractionUsecase) ToggleFavorite(ctx context.Context, userID, appID int) (bool, error) {
	isFav, err := u.favRepo.IsFavorite(ctx, userID, appID)
	if err != nil {
		return false, fmt.Errorf("failed to check favorite status: %w", err)
	}

	if isFav {
		if err := u.favRepo.Remove(ctx, userID, appID); err != nil {
			return true, fmt.Errorf("failed to remove favorite: %w", err)
		}
		return false, nil
	}

	if err := u.favRepo.Add(ctx, userID, appID); err != nil {
		return false, fmt.Errorf("failed to add favorite: %w", err)
	}
	return true, nil
}

// ListFavorites returns all favorite applications for a user
func (u *InteractionUsecase) ListFavorites(ctx context.Context, userID int) ([]domain.App, error) {
	return u.favRepo.ListByUser(ctx, userID)
}

// ListRecents returns recently opened applications for a user
func (u *InteractionUsecase) ListRecents(ctx context.Context, userID int, limit int) ([]domain.App, error) {
	return u.clickLogRepo.GetRecentByUser(ctx, userID, limit)
}
