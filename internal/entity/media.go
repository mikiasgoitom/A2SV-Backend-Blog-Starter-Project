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

// // NewMedia creates a new media instance
// func NewMedia(fileName, url, mimeType string, fileSize int64, uploadedByUserID uuid.UUID) *Media {
// 	return &Media{
// 		ID:               uuid.New(),
// 		FileName:         fileName,
// 		URL:              url,
// 		MimeType:         mimeType,
// 		FileSize:         fileSize,
// 		UploadedByUserID: uploadedByUserID,
// 		CreatedAt:        time.Now(),
// 	}
// }

// // IsImage checks if the media is an image
// func (m *Media) IsImage() bool {
// 	return m.MimeType[:5] == "image"
// }

// // IsVideo checks if the media is a video
// func (m *Media) IsVideo() bool {
// 	return m.MimeType[:5] == "video"
// }

// // GetFileSizeInMB returns the file size in megabytes
// func (m *Media) GetFileSizeInMB() float64 {
// 	return float64(m.FileSize) / (1024 * 1024)
// } 