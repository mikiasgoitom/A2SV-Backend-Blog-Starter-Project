package entity

import (
	"time"

	"github.com/google/uuid"
)

// Like represents a like on a blog post or comment
type Like struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	TargetID   uuid.UUID `json:"target_id" db:"target_id"`
	TargetType TargetType `json:"target_type" db:"target_type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// TargetType represents the type of entity being liked
type TargetType string

const (
	TargetTypeBlog    TargetType = "blog"
	TargetTypeComment TargetType = "comment"
)

// // NewLike creates a new like instance
// func NewLike(userID, targetID uuid.UUID, targetType TargetType) *Like {
// 	return &Like{
// 		ID:         uuid.New(),
// 		UserID:     userID,
// 		TargetID:   targetID,
// 		TargetType: targetType,
// 		CreatedAt:  time.Now(),
// 	}
// }

// // NewBlogLike creates a new like for a blog post
// func NewBlogLike(userID, blogID uuid.UUID) *Like {
// 	return NewLike(userID, blogID, TargetTypeBlog)
// }

// // NewCommentLike creates a new like for a comment
// func NewCommentLike(userID, commentID uuid.UUID) *Like {
// 	return NewLike(userID, commentID, TargetTypeComment)
// }

// // IsBlogLike checks if this like is for a blog post
// func (l *Like) IsBlogLike() bool {
// 	return l.TargetType == TargetTypeBlog
// }

// // IsCommentLike checks if this like is for a comment
// func (l *Like) IsCommentLike() bool {
// 	return l.TargetType == TargetTypeComment
// } 