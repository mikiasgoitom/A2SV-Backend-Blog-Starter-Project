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

// // NewTag creates a new tag instance
// func NewTag(name, slug string) *Tag {
// 	return &Tag{
// 		ID:        uuid.New(),
// 		Name:      name,
// 		Slug:      slug,
// 		CreatedAt: time.Now(),
// 	}
// }

// // UpdateTag updates the tag information
// func (t *Tag) UpdateTag(name, slug string) {
// 	t.Name = name
// 	t.Slug = slug
// } 