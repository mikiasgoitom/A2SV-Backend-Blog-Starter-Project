package contract

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// ICommentRepository defines the interface for comment data persistence.
type ICommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment) error
	GetCommentByID(ctx context.Context, commentID uuid.UUID) (*entity.Comment, error)
	GetCommentsByBlogID(ctx context.Context, blogID uuid.UUID, page, pageSize int) ([]*entity.Comment, int64, error)
	UpdateComment(ctx context.Context, commentID uuid.UUID, updates map[string]interface{}) error
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
}
