package contract

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// IBlogRepository provides methods for managing blog data in the database.
type IBlogRepository interface {
	CreateBlog(ctx context.Context, blog *entity.Blog) error
	GetBlogByID(ctx context.Context, blogID uuid.UUID) (*entity.Blog, error)
	GetBlogs(ctx context.Context, filterOptions *BlogFilterOptions) ([]*entity.Blog, int64, error)
	UpdateBlog(ctx context.Context, blogID uuid.UUID, updates map[string]interface{}) error
	DeleteBlog(ctx context.Context, blogID uuid.UUID) error
	SearchBlogs(ctx context.Context, query string, filterOptions *BlogFilterOptions) ([]*entity.Blog, int64, error)
	IncrementViewCount(ctx context.Context, blogID uuid.UUID) error
	IncrementLikeCount(ctx context.Context, blogID uuid.UUID) error
	DecrementLikeCount(ctx context.Context, blogID uuid.UUID) error
	IncrementDislikeCount(ctx context.Context, blogID uuid.UUID) error
	DecrementDislikeCount(ctx context.Context, blogID uuid.UUID) error
	IncrementCommentCount(ctx context.Context, blogID uuid.UUID) error
	DecrementCommentCount(ctx context.Context, blogID uuid.UUID) error
	GetBlogCounts(ctx context.Context, blogID uuid.UUID) (viewCount, likeCount, dislikeCount, commentCount int, err error)
	AddTagsToBlog(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error
	RemoveTagsFromBlog(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error
}

// BlogFilterOptions encapsulates filtering, pagination, and sorting parameters for blog retrieval.
type BlogFilterOptions struct {
	Page      int
	PageSize  int
	SortBy    string // e.g., "createdAt", "viewCount"
	SortOrder string // e.g., "asc", "desc"
	DateFrom  *time.Time
	DateTo    *time.Time
	MinViews  *int
	MaxViews  *int
	MinLikes  *int
	MaxLikes  *int
	AuthorID  *uuid.UUID
	TagIDs    []uuid.UUID
}
