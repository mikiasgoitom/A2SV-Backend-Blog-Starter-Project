package contract

import (
	"context"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"

	"github.com/google/uuid"
)

// IBlogRepository provides methods for managing blog data in the database.
type IBlogRepository interface {
	CreateBlog(ctx context.Context, blog *entity.Blog) error
	GetBlogs(ctx context.Context, filterOptions *BlogFilterOptions) ([]*entity.Blog, int64, error)
	GetBlogByID(ctx context.Context, blogID uuid.UUID) (*entity.Blog, error)
	UpdateBlog(ctx context.Context, blogID uuid.UUID, updates map[string]interface{}) error
	DeleteBlog(ctx context.Context, blogID uuid.UUID) error
	SearchBlogs(ctx context.Context, query string, filterOptions *BlogFilterOptions) ([]*entity.Blog, int64, error)
	IncrementViewCount(ctx context.Context, blogID uuid.UUID) error
	AddTagsToBlog(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error
	GetBlogsByTagID(ctx context.Context, tagID uuid.UUID, opts *BlogFilterOptions) ([]*entity.Blog, int64, error)
	RemoveTagsFromBlog(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error
}

// BlogFilterOptions encapsulates filtering, pagination, and sorting parameters for blog retrieval.
type BlogFilterOptions struct {
	Page         int
	PageSize     int
	SortBy       string // e.g., "created_at", "view_count"
	SortOrder    string // e.g., "asc", "desc"
	FilterByDate *time.Time
}
