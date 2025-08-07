package contract

import (
	"context"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// MediaFilterOptions holds database-agnostic parameters for filtering, sorting, and pagination.
type MediaFilterOptions struct {
	UploadedByUserID *string
	MimeType         *string
	Page             int64
	Limit            int64
	SortBy           string // e.g., "created_at", "file_name"
	SortOrder        string // "asc" or "desc"
}

// IMediaRepository defines the interface for media data persistence.
type IMediaRepository interface {
	CreateMedia(ctx context.Context, media *entity.Media) error

	GetMediaByID(ctx context.Context, mediaID string) (*entity.Media, error)
	GetMediaByBlogID(ctx context.Context, blogID string) ([]*entity.Media, error)
	GetMedia(ctx context.Context, opts *MediaFilterOptions) ([]*entity.Media, error)
	// GetMediaByURL retrieves a media record by its URL.
   //  GetMediaByUrl(ctx context.Context, url string) (*entity.Media, error)


	UpdateMedia(ctx context.Context, mediaID string, updates map[string]interface{}) error

	DeleteMedia(ctx context.Context, mediaID string) error
	//
	GetAvatarURLForUser(ctx context.Context, userID string) (string, error)

	// SetAvatarURLForUser updates the user's avatar URL.
	SetAvatarURLForUser(ctx context.Context, userID string, url string) error

	// UnsetAvatarURLForUser removes the user's avatar URL.
	UnsetAvatarURLForUser(ctx context.Context, userID string) error
}