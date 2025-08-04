package contract

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// ITagRepository defines the interface for tag data persistence.
type ITagRepository interface {
	CreateTag(ctx context.Context, tag *entity.Tag) error
	GetTagByID(ctx context.Context, tagID uuid.UUID) (*entity.Tag, error)
	GetTagByName(ctx context.Context, name string) (*entity.Tag, error)
	GetAllTags(ctx context.Context) ([]*entity.Tag, error)
	UpdateTag(ctx context.Context, tagID uuid.UUID, updates map[string]interface{}) error
	DeleteTag(ctx context.Context, tagID uuid.UUID) error
}
