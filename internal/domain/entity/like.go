package entity

import (
	"time"
)

// LikeType represents the type of reaction (like or dislike)
type LikeType string

const (
	LIKE_TYPE_LIKE    LikeType = "like"
	LIKE_TYPE_DISLIKE LikeType = "dislike"
)

// Like represents a like on a blog post or comment
type Like struct {
	ID         string     `json:"id" bson:"id"`
	UserID     string     `json:"user_id" bson:"user_id"`
	TargetID   string     `json:"target_id" bson:"target_id"`
	TargetType TargetType `json:"target_type" bson:"target_type"`
	Type       LikeType   `json:"type" bson:"type"`
	IsDeleted  bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
}

// TargetType represents the type of entity being liked
type TargetType string

const (
	TargetTypeBlog    TargetType = "blog"
	TargetTypeComment TargetType = "comment"
)
