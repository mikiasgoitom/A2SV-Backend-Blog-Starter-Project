package entity

import "time"

// BlogView represents a record of a user viewing a blog, used for tracking and analysis.
type BlogView struct {
	BlogID    string    `bson:"blog_id"`
	UserID    string    `bson:"user_id,omitempty"`
	IPAddress string    `bson:"ip_address"`
	UserAgent string    `bson:"user_agent"`
	ViewedAt  time.Time `bson:"viewed_at"`
}
