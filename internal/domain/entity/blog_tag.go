package entity

import (
	"github.com/google/uuid"
)

// BlogTag represents the many-to-many relationship between blogs and tags
type BlogTag struct {
	BlogID uuid.UUID `json:"blog_id" db:"blog_id"`
	TagID  uuid.UUID `json:"tag_id" db:"tag_id"`
}