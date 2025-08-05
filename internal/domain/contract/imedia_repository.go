package contract

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// IMediaRepository defines the interface for media data persistence.
type IMediaRepository interface {
	CreateMedia(ctx context.Context, media *entity.Media) error
	GetMediaByID(ctx context.Context, mediaID uuid.UUID) (*entity.Media, error)
	DeleteMedia(ctx context.Context, mediaID uuid.UUID) error
}
