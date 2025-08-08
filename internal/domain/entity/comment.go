package entity

import (
	"time"
)

// Comment represents a comment on a blog post
type Comment struct {
	ID        string    `json:"id" bson:"id"`
	BlogID    string    `json:"blog_id" bson:"blog_id"`
	ParentID  *string   `json:"parent_id" bson:"parent_id"`
	AuthorID  string    `json:"author_id" bson:"author_id"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	IsDeleted bool      `json:"is_deleted" bson:"is_deleted"`
}
