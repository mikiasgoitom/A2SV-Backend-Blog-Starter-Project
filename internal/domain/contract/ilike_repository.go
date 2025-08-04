package contract

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// ILikeRepository defines the interface for like data persistence.
type ILikeRepository interface {
	CreateLike(ctx context.Context, like *entity.Like) error
	DeleteLike(ctx context.Context, likeID uuid.UUID) error
	GetLikeByUserIDAndTargetID(ctx context.Context, userID, targetID uuid.UUID) (*entity.Like, error)
	CountLikesByTargetID(ctx context.Context, targetID uuid.UUID) (int64, error)
}
