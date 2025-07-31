package entity

import (
	"time"

	"github.com/google/uuid"
)

// Blog represents a blog post in the system
type Blog struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	Content         string     `json:"content" db:"content"`
	AuthorID        uuid.UUID  `json:"author_id" db:"author_id"`
	Slug            string     `json:"slug" db:"slug"`
	Status          BlogStatus `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	PublishedAt     *time.Time `json:"published_at" db:"published_at"`
	ViewCount       int        `json:"view_count" db:"view_count"`
	FeaturedImageID *uuid.UUID `json:"featured_image_id" db:"featured_image_id"`
	IsDeleted       bool       `json:"is_deleted" db:"is_deleted"`
}

// BlogStatus represents the status of a blog post
type BlogStatus string

const (
	BlogStatusDraft     BlogStatus = "draft"
	BlogStatusPublished BlogStatus = "published"
	BlogStatusArchived  BlogStatus = "archived"
)

// // NewBlog creates a new blog instance
// func NewBlog(title, content string, authorID uuid.UUID, slug string) *Blog {
// 	now := time.Now()
// 	return &Blog{
// 		ID:        uuid.New(),
// 		Title:     title,
// 		Content:   content,
// 		AuthorID:  authorID,
// 		Slug:      slug,
// 		Status:    BlogStatusDraft,
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 		ViewCount: 0,
// 		IsDeleted: false,
// 	}
// }

// // Publish publishes the blog post
// func (b *Blog) Publish() {
// 	now := time.Now()
// 	b.Status = BlogStatusPublished
// 	b.PublishedAt = &now
// 	b.UpdatedAt = now
// }

// // Archive archives the blog post
// func (b *Blog) Archive() {
// 	b.Status = BlogStatusArchived
// 	b.UpdatedAt = time.Now()
// }

// // IncrementViewCount increments the view count
// func (b *Blog) IncrementViewCount() {
// 	b.ViewCount++
// 	b.UpdatedAt = time.Now()
// }

// // SetFeaturedImage sets the featured image for the blog
// func (b *Blog) SetFeaturedImage(imageID uuid.UUID) {
// 	b.FeaturedImageID = &imageID
// 	b.UpdatedAt = time.Now()
// }

// // SoftDelete marks the blog as deleted without removing it from the database
// func (b *Blog) SoftDelete() {
// 	b.IsDeleted = true
// 	b.UpdatedAt = time.Now()
// }

// // IsPublished checks if the blog is published
// func (b *Blog) IsPublished() bool {
// 	return b.Status == BlogStatusPublished
// }

// // IsDraft checks if the blog is in draft status
// func (b *Blog) IsDraft() bool {
// 	return b.Status == BlogStatusDraft
// }

// // IsArchived checks if the blog is archived
// func (b *Blog) IsArchived() bool {
// 	return b.Status == BlogStatusArchived
// }