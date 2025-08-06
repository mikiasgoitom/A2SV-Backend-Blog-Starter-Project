package contract

import (
	"context"
)

// BlogPopularityRepository extends the BlogRepository with popularity tracking methods
type BlogPopularityRepository interface {
	// Like/Dislike methods
	IncrementLikeCount(ctx context.Context, blogID string) error
	DecrementLikeCount(ctx context.Context, blogID string) error
	IncrementDislikeCount(ctx context.Context, blogID string) error
	DecrementDislikeCount(ctx context.Context, blogID string) error

	// Comment methods
	IncrementCommentCount(ctx context.Context, blogID string) error
	DecrementCommentCount(ctx context.Context, blogID string) error

	// User action tracking
	AddUserLike(ctx context.Context, blogID, userID string, likeType string) error
	RemoveUserLike(ctx context.Context, blogID, userID string) error
	HasUserLiked(ctx context.Context, blogID, userID string) (string, bool, error)

	// Get current counts
	GetBlogCounts(ctx context.Context, blogID string) (viewCount, likeCount, dislikeCount, commentCount int, err error)
}
