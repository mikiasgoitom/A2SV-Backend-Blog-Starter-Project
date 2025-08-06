package contract

import (
	"context"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// ILikeRepository defines the interface for like data persistence.
type ILikeRepository interface {
	CreateLike(ctx context.Context, like *entity.Like) error
	DeleteLike(ctx context.Context, likeID string) error
	GetLikeByUserIDAndTargetID(ctx context.Context, userID, targetID string) (*entity.Like, error)
	CountLikesByTargetID(ctx context.Context, targetID string) (int64, error)
}
