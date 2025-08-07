package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
)

type InteractionAction string

const (
	InteractionActionLike   InteractionAction = "like"
	InteractionActionUnlike InteractionAction = "unlike"
)

type IInteractionUseCase interface {
	LikeBlog(ctx context.Context, blogID, userID string) error
	UnlikeBlog(ctx context.Context, blogID, userID string) error
}

type InteractionUseCaseImpl struct {
	blogRepo contract.IBlogRepository
	logger   AppLogger
}

func NewInteractionUseCase(blogRepo contract.IBlogRepository, logger AppLogger) *InteractionUseCaseImpl {
	return &InteractionUseCaseImpl{
		blogRepo: blogRepo,
		logger:   logger,
	}
}

func (uc *InteractionUseCaseImpl) LikeBlog(ctx context.Context, blogID, userID string) error {
	if blogID == "" {
		return errors.New("blog ID is required")
	}
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Check if the user has already liked the blog
	likeType, hasLiked, err := uc.blogRepo.HasUserLiked(ctx, blogID, userID)
	if err != nil {
		uc.logger.Errorf("failed to check if user has liked blog: %v", err)
		return fmt.Errorf("failed to check if user has liked blog: %w", err)
	}

	if hasLiked && likeType == "like" {
		return errors.New("user has already liked this blog")
	}

	// Add the like
	if err := uc.blogRepo.AddUserLike(ctx, blogID, userID, "like"); err != nil {
		uc.logger.Errorf("failed to add user like: %v", err)
		return fmt.Errorf("failed to add user like: %w", err)
	}

	// Increment the like count
	if err := uc.blogRepo.IncrementLikeCount(ctx, blogID); err != nil {
		uc.logger.Errorf("failed to increment like count: %v", err)
		return fmt.Errorf("failed to increment like count: %w", err)
	}

	uc.logger.Infof("User %s liked blog %s", userID, blogID)
	return nil
}

func (uc *InteractionUseCaseImpl) UnlikeBlog(ctx context.Context, blogID, userID string) error {
	if blogID == "" {
		return errors.New("blog ID is required")
	}
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Check if the user has liked the blog
	likeType, hasLiked, err := uc.blogRepo.HasUserLiked(ctx, blogID, userID)
	if err != nil {
		uc.logger.Errorf("failed to check if user has liked blog: %v", err)
		return fmt.Errorf("failed to check if user has liked blog: %w", err)
	}

	if !hasLiked || likeType != "like" {
		return errors.New("user has not liked this blog")
	}

	// Remove the like
	if err := uc.blogRepo.RemoveUserLike(ctx, blogID, userID); err != nil {
		uc.logger.Errorf("failed to remove user like: %v", err)
		return fmt.Errorf("failed to remove user like: %w", err)
	}

	// Decrement the like count
	if err := uc.blogRepo.DecrementLikeCount(ctx, blogID); err != nil {
		uc.logger.Errorf("failed to decrement like count: %v", err)
		return fmt.Errorf("failed to decrement like count: %w", err)
	}

	uc.logger.Infof("User %s unliked blog %s", userID, blogID)
	return nil
}
