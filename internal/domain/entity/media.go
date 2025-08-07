package entity

import (
	"time"
)

// Media represents an uploaded media file
type Media struct {
	ID               string    `json:"id" db:"id"`
	FileName         string    `json:"file_name" db:"file_name"`
	URL              string    `json:"url" db:"url"`
	MimeType         string    `json:"mime_type" db:"mime_type"`
	FileSize         int64     `json:"file_size" db:"file_size"`
	UploadedByUserID string    `json:"uploaded_by_user_id" db:"uploaded_by_user_id"`
	BlogID           string    `json:"blog_id,omitempty" db:"blog_id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	IsDeleted        bool      `json:"is_deleted,omitempty" db:"is_deleted"`
}
