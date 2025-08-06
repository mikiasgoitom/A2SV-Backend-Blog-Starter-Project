package contract

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// ILikeRepository defines the interface for reaction data persistence.
type ILikeRepository interface {
	CreateReaction(ctx context.Context, like *entity.Like) error
	DeleteReaction(ctx context.Context, reactionID uuid.UUID) error
	GetReactionByUserIDAndTargetID(ctx context.Context, userID, targetID uuid.UUID) (*entity.Like, error)
	GetReactionByUserIDTargetIDAndType(ctx context.Context, userID, targetID uuid.UUID, reactionType entity.LikeType) (*entity.Like, error)
	CountLikesByTargetID(ctx context.Context, targetID uuid.UUID) (int64, error)
	CountDislikesByTargetID(ctx context.Context, targetID uuid.UUID) (int64, error)
}
