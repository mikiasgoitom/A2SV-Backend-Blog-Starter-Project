package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	mongodb "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/repository/mongodb"
)

// LikeUsecase handles the business logic for managing likes and dislikes.
type LikeUsecase struct {
	likeRepo contract.ILikeRepository
	blogRepo contract.IBlogRepository // Add blogRepo for updating popularity
}

// NewLikeUsecase creates and returns a new LikeUsecase instance.
func NewLikeUsecase(likeRepo contract.ILikeRepository, blogRepo contract.IBlogRepository) *LikeUsecase {
	return &LikeUsecase{
		likeRepo: likeRepo,
		blogRepo: blogRepo,
	}
}

// ToggleLike handles the logic for liking and unliking a target.
func (u *LikeUsecase) ToggleLike(ctx context.Context, userID, targetID string, targetType entity.TargetType) error {
	existingReaction, err := u.likeRepo.GetReactionByUserIDAndTargetID(ctx, userID, targetID)
	if err != nil {
		if errors.Is(err, mongodb.ErrReactionNotFound) {
			existingReaction = nil
		} else {
			return fmt.Errorf("failed to retrieve existing reaction: %w", err)
		}
	}

	var resultErr error
	if existingReaction != nil {
		if existingReaction.Type == entity.LIKE_TYPE_LIKE {
			// User is unliking a target they've already liked.
			resultErr = u.likeRepo.DeleteReaction(ctx, existingReaction.ID)
		} else {
			// User is changing a 'dislike' to a 'like'.
			existingReaction.Type = entity.LIKE_TYPE_LIKE
			resultErr = u.likeRepo.CreateReaction(ctx, existingReaction)
		}
	} else {
		// No reaction exists, create a new one.
		newLike := &entity.Like{
			UserID:     userID,
			TargetID:   targetID,
			TargetType: targetType,
			Type:       entity.LIKE_TYPE_LIKE,
		}
		resultErr = u.likeRepo.CreateReaction(ctx, newLike)
	}

	// Update blog like_count and popularity if this is a blog like/dislike
	if targetType == entity.TargetTypeBlog && u.blogRepo != nil {
		// Recalculate like/dislike counts from likes collection
		likes, err1 := u.likeRepo.CountLikesByTargetID(ctx, targetID)
		dislikes, err2 := u.likeRepo.CountDislikesByTargetID(ctx, targetID)
		if err1 == nil && err2 == nil {
			// Get current view and comment counts from blog
			blog, err := u.blogRepo.GetBlogByID(ctx, targetID)
			views := 0
			comments := 0
			if err == nil && blog != nil {
				views = blog.ViewCount
				comments = blog.CommentCount
			}
			popularity := calculatePopularity(views, int(likes), int(dislikes), comments)
			updates := map[string]interface{}{
				"like_count":    likes,
				"dislike_count": dislikes,
				"popularity":    popularity,
			}
			_ = u.blogRepo.UpdateBlog(ctx, targetID, updates)
		}
	}
	return resultErr
}

// ToggleDislike handles the logic for disliking and undisliking a target.
func (u *LikeUsecase) ToggleDislike(ctx context.Context, userID, targetID string, targetType entity.TargetType) error {
	existingReaction, err := u.likeRepo.GetReactionByUserIDAndTargetID(ctx, userID, targetID)
	if err != nil {
		if errors.Is(err, mongodb.ErrReactionNotFound) {
			existingReaction = nil
		} else {
			return fmt.Errorf("failed to retrieve existing reaction: %w", err)
		}
	}

	var resultErr error
	if existingReaction != nil {
		if existingReaction.Type == entity.LIKE_TYPE_DISLIKE {
			// User is undisliking a target they've already disliked.
			resultErr = u.likeRepo.DeleteReaction(ctx, existingReaction.ID)
		} else {
			// User is changing a 'like' to a 'dislike'.
			existingReaction.Type = entity.LIKE_TYPE_DISLIKE
			resultErr = u.likeRepo.CreateReaction(ctx, existingReaction)
		}
	} else {
		// No reaction exists, create a new one.
		newDislike := &entity.Like{
			UserID:     userID,
			TargetID:   targetID,
			TargetType: targetType,
			Type:       entity.LIKE_TYPE_DISLIKE,
		}
		resultErr = u.likeRepo.CreateReaction(ctx, newDislike)
	}

	// Update blog dislike_count and popularity if this is a blog like/dislike
	if targetType == entity.TargetTypeBlog && u.blogRepo != nil {
		// Recalculate like/dislike counts from likes collection
		likes, err1 := u.likeRepo.CountLikesByTargetID(ctx, targetID)
		dislikes, err2 := u.likeRepo.CountDislikesByTargetID(ctx, targetID)
		if err1 == nil && err2 == nil {
			// Get current view and comment counts from blog
			blog, err := u.blogRepo.GetBlogByID(ctx, targetID)
			views := 0
			comments := 0
			if err == nil && blog != nil {
				views = blog.ViewCount
				comments = blog.CommentCount
			}
			popularity := calculatePopularity(views, int(likes), int(dislikes), comments)
			updates := map[string]interface{}{
				"like_count":    likes,
				"dislike_count": dislikes,
				"popularity":    popularity,
			}
			_ = u.blogRepo.UpdateBlog(ctx, targetID, updates)
		}
	}
	return resultErr
}

// GetUserReaction retrieves the active reaction (if any) a user has on a specific target.
func (u *LikeUsecase) GetUserReaction(ctx context.Context, userID, targetID string) (*entity.Like, error) {
	like, err := u.likeRepo.GetReactionByUserIDAndTargetID(ctx, userID, targetID)
	if err != nil {
		if errors.Is(err, mongodb.ErrReactionNotFound) {
			// The use case should handle this specific error and return nil, nil
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user's reaction: %w", err)
	}
	return like, nil
}

// GetReactionCounts retrieves the total number of likes and dislikes for a specific target.
func (u *LikeUsecase) GetReactionCounts(ctx context.Context, targetID string) (likes, dislikes int64, err error) {
	likes, err = u.likeRepo.CountLikesByTargetID(ctx, targetID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count likes for target %s: %w", targetID, err)
	}

	dislikes, err = u.likeRepo.CountDislikesByTargetID(ctx, targetID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count dislikes for target %s: %w", targetID, err)
	}

	return likes, dislikes, nil
}
