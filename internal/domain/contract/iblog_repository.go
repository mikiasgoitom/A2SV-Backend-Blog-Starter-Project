package contract

import (
	"context"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// IBlogRepository provides methods for managing blog data in the database.
type IBlogRepository interface {
	CreateBlog(ctx context.Context, blog *entity.Blog) error
	GetBlogs(ctx context.Context, filterOptions *BlogFilterOptions) ([]*entity.Blog, int64, error)
	GetBlogByID(ctx context.Context, blogID string) (*entity.Blog, error)
	UpdateBlog(ctx context.Context, blogID string, updates map[string]interface{}) error
	DeleteBlog(ctx context.Context, blogID string) error
	SearchBlogs(ctx context.Context, query string, filterOptions *BlogFilterOptions) ([]*entity.Blog, int64, error)
	IncrementViewCount(ctx context.Context, blogID string) error
	AddTagsToBlog(ctx context.Context, blogID string, tagIDs []string) error
	GetBlogsByTagID(ctx context.Context, tagID string, opts *BlogFilterOptions) ([]*entity.Blog, int64, error)
	GetBlogsByTagIDs(ctx context.Context, tagIDs []string, page int, pageSize int) ([]*entity.Blog, int64, error)
	RemoveTagsFromBlog(ctx context.Context, blogID string, tagIDs []string) error
}

// BlogFilterOptions encapsulates filtering, pagination, and sorting parameters for blog retrieval.
type BlogFilterOptions struct {
	Page      int
	PageSize  int
	SortBy    string // e.g., "created_at", "view_count"
	SortOrder string // e.g., "asc", "desc"
	DateFrom  *time.Time
	DateTo    *time.Time
	MinViews  *int
	MaxViews  *int
	MinLikes  *int
	MaxLikes  *int
	AuthorID  *string
}
