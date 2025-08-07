package entity

import (
	"time"
)

// Tag represents a tag for categorizing blog posts
type Tag struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Slug      string    `json:"slug" bson:"slug"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
