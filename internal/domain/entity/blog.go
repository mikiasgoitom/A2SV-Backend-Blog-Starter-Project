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