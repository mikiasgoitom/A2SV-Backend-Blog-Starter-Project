package entity

import (
	"time"

	"github.com/google/uuid"
)

// LikeType represents the type of reaction (like or dislike)
type LikeType string

const (
	LIKE_TYPE_LIKE    LikeType = "like"
	LIKE_TYPE_DISLIKE LikeType = "dislike"
)

// Like represents a like on a blog post or comment
type Like struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	TargetID   uuid.UUID  `json:"target_id" db:"target_id"`
	TargetType TargetType `json:"target_type" db:"target_type"`
	Type       LikeType   `json:"type" db:"type"`
	IsDeleted  bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// TargetType represents the type of entity being liked
type TargetType string

const (
	TargetTypeBlog    TargetType = "blog"
	TargetTypeComment TargetType = "comment"
)
