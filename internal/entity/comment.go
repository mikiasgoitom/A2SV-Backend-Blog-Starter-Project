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

// // NewComment creates a new comment instance
// func NewComment(blogID, authorID uuid.UUID, content string, parentID *uuid.UUID) *Comment {
// 	now := time.Now()
// 	return &Comment{
// 		ID:        uuid.New(),
// 		BlogID:    blogID,
// 		ParentID:  parentID,
// 		AuthorID:  authorID,
// 		Content:   content,
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 		IsDeleted: false,
// 	}
// }

// // UpdateContent updates the comment content
// func (c *Comment) UpdateContent(content string) {
// 	c.Content = content
// 	c.UpdatedAt = time.Now()
// }

// // SoftDelete marks the comment as deleted without removing it from the database
// func (c *Comment) SoftDelete() {
// 	c.IsDeleted = true
// 	c.UpdatedAt = time.Now()
// }

// // IsReply checks if this comment is a reply to another comment
// func (c *Comment) IsReply() bool {
// 	return c.ParentID != nil
// }

// // IsTopLevel checks if this comment is a top-level comment (not a reply)
// func (c *Comment) IsTopLevel() bool {
// 	return c.ParentID == nil
// } 