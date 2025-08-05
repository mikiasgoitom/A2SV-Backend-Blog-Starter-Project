package contract

import (
	"context"

	"github.com/google/uuid"
)

// BlogPopularityRepository extends the BlogRepository with popularity tracking methods
type BlogPopularityRepository interface {
	// Like/Dislike methods
	IncrementLikeCount(ctx context.Context, blogID uuid.UUID) error
	DecrementLikeCount(ctx context.Context, blogID uuid.UUID) error
	IncrementDislikeCount(ctx context.Context, blogID uuid.UUID) error
	DecrementDislikeCount(ctx context.Context, blogID uuid.UUID) error

	// Comment methods
	IncrementCommentCount(ctx context.Context, blogID uuid.UUID) error
	DecrementCommentCount(ctx context.Context, blogID uuid.UUID) error

	// User action tracking
	AddUserLike(ctx context.Context, blogID, userID uuid.UUID, likeType string) error
	RemoveUserLike(ctx context.Context, blogID, userID uuid.UUID) error
	HasUserLiked(ctx context.Context, blogID, userID uuid.UUID) (string, bool, error)

	// Get current counts
	GetBlogCounts(ctx context.Context, blogID uuid.UUID) (viewCount, likeCount, dislikeCount, commentCount int, err error)
}
