package usecasecontract

import (
	"context"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

type ILikeUseCase interface {
	ToggleLike(ctx context.Context, userID, targetID string, targetType entity.TargetType) error
	ToggleDislike(ctx context.Context, userID, targetID string, targetType entity.TargetType) error
	GetUserReaction(ctx context.Context, userID, targetID string) (*entity.Like, error)
	GetReactionCounts(ctx context.Context, targetID string) (likes, dislikes int64, err error)
}
