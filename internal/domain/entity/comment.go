package entity

import (
	"time"
)

// Comment represents a comment on a blog post with advanced reply-to-reply support
type Comment struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	BlogID         string    `json:"blog_id" bson:"blog_id"`
	Type           string    `json:"type" bson:"type"` // "comment" or "reply"
	ParentID       *string   `json:"parent_id" bson:"parent_id"`
	TargetID       *string   `json:"target_id" bson:"target_id"`
	AuthorID       string    `json:"author_id" bson:"author_id"`
	AuthorName     string    `json:"author_name" bson:"author_name"`
	TargetUserName string    `json:"target_user_name" bson:"target_user_name"`
	Content        string    `json:"content" bson:"content"`
	Status         string    `json:"status" bson:"status"`
	LikeCount      int       `json:"like_count" bson:"like_count"`
	ReplyCount     int       `json:"reply_count" bson:"reply_count"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
}

// CommentThread represents a comment with its nested replies
type CommentThread struct {
	Comment *Comment         `json:"comment"`
	Replies []*CommentThread `json:"replies"`
	Depth   int              `json:"depth"`
}

// CommentLike represents a user's like on a comment
type CommentLike struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	CommentID string    `json:"comment_id" bson:"comment_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// CommentReport represents a report against a comment
type CommentReport struct {
	ID         string     `json:"id" bson:"_id,omitempty"`
	CommentID  string     `json:"comment_id" bson:"comment_id"`
	ReporterID string     `json:"reporter_id" bson:"reporter_id"`
	Reason     string     `json:"reason" bson:"reason"`
	Details    string     `json:"details" bson:"details"`
	Status     string     `json:"status" bson:"status"` // pending, reviewed, dismissed
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	ReviewedAt *time.Time `json:"reviewed_at" bson:"reviewed_at"`
	ReviewedBy *string    `json:"reviewed_by" bson:"reviewed_by"`
}
