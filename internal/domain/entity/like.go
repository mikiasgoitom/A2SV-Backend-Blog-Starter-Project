package entity

import (
	"time"
)

// Like represents a like on a blog post or comment
type Like struct {
	ID         string     `json:"id" bson:"id"`
	UserID     string     `json:"user_id" bson:"user_id"`
	TargetID   string     `json:"target_id" bson:"target_id"`
	TargetType TargetType `json:"target_type" bson:"target_type"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
}

// TargetType represents the type of entity being liked
type TargetType string

const (
	TargetTypeBlog    TargetType = "blog"
	TargetTypeComment TargetType = "comment"
)
