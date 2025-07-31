package entity

import (
	"time"

	"github.com/google/uuid"
)

// Media represents an uploaded media file
type Media struct {
	ID                uuid.UUID `json:"id" db:"id"`
	FileName          string    `json:"file_name" db:"file_name"`
	URL               string    `json:"url" db:"url"`
	MimeType          string    `json:"mime_type" db:"mime_type"`
	FileSize          int64     `json:"file_size" db:"file_size"`
	UploadedByUserID  uuid.UUID `json:"uploaded_by_user_id" db:"uploaded_by_user_id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}