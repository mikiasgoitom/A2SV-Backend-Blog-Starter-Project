package contract

import (
	"context"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// IUserAvatarRepository defines the interface for user avatar data persistence.
// This interface separates business logic from data storage details.
type IUserAvatarRepository interface {
	// CreateAvatarURI saves a new avatar URI for a user.
	CreateAvatarURI(ctx context.Context, userID, uri string) error

	// ReadAvatarURI retrieves the avatar URI from a user's record.
	ReadAvatarURI(ctx context.Context, userID string) (string, error)

	// UpdateAvatarURI updates the user's avatar URI with a new one.
	UpdateAvatarURI(ctx context.Context, userID, newURI string) error

	// DeleteAvatarURI removes the avatar URI from a user's record.
	DeleteAvatarURI(ctx context.Context, userID string) error

	// ReadMediaByURI retrieves media metadata using the avatar's URI.
	ReadMediaByURI(ctx context.Context, uri string) (*entity.Media, error)
}
