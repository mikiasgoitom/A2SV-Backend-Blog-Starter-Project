package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// LikeUsecase handles the business logic for managing likes and dislikes.
type LikeUsecase struct {
	likeRepo contract.ILikeRepository
}

// NewLikeUsecase creates and returns a new LikeUsecase instance.
func NewLikeUsecase(likeRepo contract.ILikeRepository) *LikeUsecase {
	return &LikeUsecase{
		likeRepo: likeRepo,
	}
}

// ToggleLike handles the logic for liking and unliking a target.
func (u *LikeUsecase) ToggleLike(ctx context.Context, userID, targetID string, targetType entity.TargetType) error {
	existingReaction, err := u.likeRepo.GetReactionByUserIDAndTargetID(ctx, userID, targetID)
	if err != nil && !errors.Is(err, errors.New("reaction not found")) {
		return fmt.Errorf("failed to retrieve existing reaction: %w", err)
	}

	if existingReaction != nil {
		if existingReaction.Type == entity.LIKE_TYPE_LIKE {
			// User is unliking a target they've already liked.
			return u.likeRepo.DeleteReaction(ctx, existingReaction.ID)
		}

		// User is changing a 'dislike' to a 'like'.
		existingReaction.Type = entity.LIKE_TYPE_LIKE
		return u.likeRepo.CreateReaction(ctx, existingReaction)
	}

	// No reaction exists, create a new one.
	newLike := &entity.Like{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
		Type:       entity.LIKE_TYPE_LIKE,
	}
	return u.likeRepo.CreateReaction(ctx, newLike)
}

// ToggleDislike handles the logic for disliking and undisliking a target.
func (u *LikeUsecase) ToggleDislike(ctx context.Context, userID, targetID string, targetType entity.TargetType) error {
	existingReaction, err := u.likeRepo.GetReactionByUserIDAndTargetID(ctx, userID, targetID)
	if err != nil && !errors.Is(err, errors.New("reaction not found")) {
		return fmt.Errorf("failed to retrieve existing reaction: %w", err)
	}

	if existingReaction != nil {
		if existingReaction.Type == entity.LIKE_TYPE_DISLIKE {
			// User is undisliking a target they've already disliked.
			return u.likeRepo.DeleteReaction(ctx, existingReaction.ID)
		}

		// User is changing a 'like' to a 'dislike'.
		existingReaction.Type = entity.LIKE_TYPE_DISLIKE
		return u.likeRepo.CreateReaction(ctx, existingReaction)
	}

	// No reaction exists, create a new one.
	newDislike := &entity.Like{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
		Type:       entity.LIKE_TYPE_DISLIKE,
	}
	return u.likeRepo.CreateReaction(ctx, newDislike)
}

// GetUserReaction retrieves the active reaction (if any) a user has on a specific target.
func (u *LikeUsecase) GetUserReaction(ctx context.Context, userID, targetID string) (*entity.Like, error) {
	like, err := u.likeRepo.GetReactionByUserIDAndTargetID(ctx, userID, targetID)
	if err != nil {
		if errors.Is(err, errors.New("reaction not found")) {
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
