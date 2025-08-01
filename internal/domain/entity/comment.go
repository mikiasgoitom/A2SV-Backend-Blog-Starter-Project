package entity

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment on a blog post
type Comment struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	BlogID    uuid.UUID  `json:"blog_id" db:"blog_id"`
	ParentID  *uuid.UUID `json:"parent_id" db:"parent_id"`
	AuthorID  uuid.UUID  `json:"author_id" db:"author_id"`
	Content   string     `json:"content" db:"content"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	IsDeleted bool       `json:"is_deleted" db:"is_deleted"`
}