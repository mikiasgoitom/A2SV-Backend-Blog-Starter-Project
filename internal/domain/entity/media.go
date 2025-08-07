package entity

import (
	"time"
)

// Media represents an uploaded media file
type Media struct {
	ID               string    `json:"id" bson:"id"`
	FileName         string    `json:"file_name" bson:"file_name"`
	URL              string    `json:"url" bson:"url"`
	MimeType         string    `json:"mime_type" bson:"mime_type"`
	FileSize         int64     `json:"file_size" bson:"file_size"`
	UploadedByUserID string    `json:"uploaded_by_user_id" bson:"uploaded_by_user_id"`
	BlogID           string    `json:"blog_id,omitempty" bson:"blog_id"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	IsDeleted        bool      `json:"is_deleted,omitempty" bson:"is_deleted"`
}
