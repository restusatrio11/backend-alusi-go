package usecase

import (
	"context"

	"backend-alusi-go/internal/repository/postgres"
	"backend-alusi-go/pkg/ai"
	"backend-alusi-go/pkg/jwt"
)

type AIUsecase struct {
	aiService *ai.AssistantService
	userRepo  *postgres.UserRepo
}

func NewAIUsecase(aiService *ai.AssistantService, userRepo *postgres.UserRepo) *AIUsecase {
	return &AIUsecase{
		aiService: aiService,
		userRepo:  userRepo,
	}
}

func (u *AIUsecase) Ask(
	ctx context.Context,
	question string,
	claims *jwt.SessionClaims,
) (*ai.AIResponse, error) {
	var userID *int
	var roleIDs []int

	if claims != nil {
		userID = &claims.UserID
		if u.userRepo != nil {
			user, err := u.userRepo.GetByID(ctx, claims.UserID)
			if err == nil && user != nil {
				for _, r := range user.Roles {
					roleIDs = append(roleIDs, r.ID)
				}
			}
		}
	}

	isPublicOnly := claims == nil

	return u.aiService.Ask(ctx, question, userID, roleIDs, isPublicOnly)
}
