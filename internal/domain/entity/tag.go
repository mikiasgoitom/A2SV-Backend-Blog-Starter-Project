package entity

import (
	"time"

	"github.com/google/uuid"
)

// Tag represents a tag for categorizing blog posts
type Tag struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}