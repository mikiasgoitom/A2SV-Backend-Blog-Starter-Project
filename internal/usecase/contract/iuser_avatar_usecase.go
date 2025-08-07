package usecasecontract

import (
	"context"
	"mime/multipart"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// IUserAvatarUseCase defines the business logic interface for managing user avatars.
// This interface orchestrates the process of handling avatars from receiving files to updating the database.
type IUserAvatarUseCase interface {
	// CreateUserAvatar handles the upload and storage of a new avatar.
	// This combines both file storage and creating the URI link in the database.
	CreateUserAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (*entity.Media, error)

	// ReadUserAvatarMetadata fetches the metadata for a user's current avatar.
	ReadUserAvatarMetadata(ctx context.Context, userID string) (*entity.Media, error)

	// UpdateUserAvatar handles the upload of a new avatar and replaces the old one.
	UpdateUserAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (*entity.Media, error)

	// DeleteUserAvatar removes the link to the avatar from the user's record and potentially the file.
	DeleteUserAvatar(ctx context.Context, userID string) error
}
