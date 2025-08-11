package entity

import (
	"time"
)

// Blog represents a blog post in the system
type Blog struct {
	ID              string     `json:"id" bson:"_id"`
	Title           string     `json:"title" bson:"title"`
	Content         string     `json:"content" bson:"content"`
	AuthorID        string     `json:"author_id" bson:"author_id"`
	Slug            string     `json:"slug" bson:"slug"`
	Status          BlogStatus `json:"status" bson:"status"`
	Tags            []string   `json:"tags" bson:"tags"`
	CreatedAt       time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" bson:"updated_at"`
	PublishedAt     *time.Time `json:"published_at" bson:"published_at"`
	ViewCount       int        `json:"view_count" bson:"view_count"`
	LikeCount       int        `json:"like_count" bson:"like_count"`
	DislikeCount    int        `json:"dislike_count" bson:"dislike_count"`
	CommentCount    int        `json:"comment_count" bson:"comment_count"`
	Popularity      float64    `json:"popularity" bson:"popularity"`
	FeaturedImageID *string    `json:"featured_image_id" bson:"featured_image_id"`
	IsDeleted       bool       `json:"is_deleted" bson:"is_deleted"`
}

// BlogStatus represents the status of a blog post
type BlogStatus string

const (
	BlogStatusDraft     BlogStatus = "draft"
	BlogStatusPublished BlogStatus = "published"
	BlogStatusArchived  BlogStatus = "archived"
)
