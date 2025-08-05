package entity

import (
	"time"

)

// Like represents a like on a blog post or comment
type Like struct {
	ID         string     `json:"id" db:"id"`
	UserID     string     `json:"user_id" db:"user_id"`
	TargetID   string     `json:"target_id" db:"target_id"`
	TargetType TargetType `json:"target_type" db:"target_type"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// TargetType represents the type of entity being liked
type TargetType string

const (
	TargetTypeBlog    TargetType = "blog"
	TargetTypeComment TargetType = "comment"
)
