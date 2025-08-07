package contract

import (
	"context"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// IBlogRepository provides methods for managing blog data in the database.
type IBlogRepository interface {
	CreateBlog(ctx context.Context, blog *entity.Blog) error
	GetBlogByID(ctx context.Context, blogID string) (*entity.Blog, error)
	GetBlogs(ctx context.Context, filterOptions *BlogFilterOptions) ([]*entity.Blog, int64, error)
	UpdateBlog(ctx context.Context, blogID string, updates map[string]interface{}) error
	DeleteBlog(ctx context.Context, blogID string) error
	SearchBlogs(ctx context.Context, query string, filterOptions *BlogFilterOptions) ([]*entity.Blog, int64, error)
	IncrementViewCount(ctx context.Context, blogID string) error
	IncrementLikeCount(ctx context.Context, blogID string) error
	DecrementLikeCount(ctx context.Context, blogID string) error
	IncrementDislikeCount(ctx context.Context, blogID string) error
	DecrementDislikeCount(ctx context.Context, blogID string) error
	IncrementCommentCount(ctx context.Context, blogID string) error
	DecrementCommentCount(ctx context.Context, blogID string) error
	GetBlogCounts(ctx context.Context, blogID string) (viewCount, likeCount, dislikeCount, commentCount int, err error)
	AddTagsToBlog(ctx context.Context, blogID string, tagIDs []string) error
	RemoveTagsFromBlog(ctx context.Context, blogID string, tagIDs []string) error
	GetBlogsByTagIDs(context.Context, []string, int, int) ([]*entity.Blog, int64, error)
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
	AuthorID  *string
	TagIDs    []string
}
