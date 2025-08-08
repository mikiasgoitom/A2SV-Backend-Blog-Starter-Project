package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "github.com/acme/avatar/config"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"

	// "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/config"

	// "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/config"
	usecasecontract "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase/contract"
)

const (
	avatarStoragePath = "uploads/avatar_uri"
	maxFileSize       = 5 * 1024 * 1024 // 5MB
)

// userAvatarUseCase implements IUserAvatarUseCase for managing user avatars
type userAvatarUseCase struct {
	userAvatarRepo contract.IUserAvatarRepository
}

// NewUserAvatarUseCase creates a new instance of IUserAvatarUseCase
func NewUserAvatarUseCase(repo contract.IUserAvatarRepository) usecasecontract.IUserAvatarUseCase {
	return &userAvatarUseCase{
		userAvatarRepo: repo,
	}
}

// CreateUserAvatar handles the upload and storage of a new avatar
func (u *userAvatarUseCase) CreateUserAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (*entity.Media, error) {
	if file.Size > maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of 5MB")
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("unsupported file type: %s", contentType)
	}

	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(avatarStoragePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Generate filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".jpg" // default extension
	}
	filename := fmt.Sprintf("%s%s", userID, ext)
	filePath := filepath.Join(avatarStoragePath, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Determine full URL for avatar
	// Explicitly set baseURL to a valid URL
	baseURL := os.Getenv("UPLOAD_BASE_URL")
	log.Printf("DEBUG: Assigned baseURL: %s", baseURL)

	// Unconditionally trim trailing slash from baseURL
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Ensure avatarPath starts with a single slash
	avatarPath := fmt.Sprintf("/uploads/avatar_uri/%s", filename)
	if !strings.HasPrefix(avatarPath, "/") {
		avatarPath = "/" + avatarPath
	}

	avatarURI := fmt.Sprintf("%s%s", baseURL, avatarPath)
	log.Printf("DEBUG: Constructed avatarURI: %s", avatarURI)
	// Save avatar URI to repository
	if err := u.userAvatarRepo.CreateAvatarURI(ctx, userID, avatarURI); err != nil {
		// Clean up file if database update fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save avatar URI: %w", err)
	}

	// Create media entity
	media := &entity.Media{
		ID:               filename,
		FileName:         file.Filename,
		URL:              avatarURI,
		MimeType:         contentType,
		FileSize:         file.Size,
		UploadedByUserID: userID,
		CreatedAt:        time.Now(),
	}

	return media, nil
}

// ReadUserAvatarMetadata fetches the metadata for a user's current avatar
func (u *userAvatarUseCase) ReadUserAvatarMetadata(ctx context.Context, userID string) (*entity.Media, error) {
	// Get the avatar URI from the repository
	uri, err := u.userAvatarRepo.ReadAvatarURI(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read avatar metadata: %w", err)
	}

	if uri == "" {
		return nil, fmt.Errorf("no avatar found for user %s", userID)
	}

	// Get file info
	filePath := strings.TrimPrefix(uri, "/")
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Create media entity
	return &entity.Media{
		ID:               filepath.Base(uri),
		FileName:         filepath.Base(uri),
		URL:              uri,
		MimeType:         "image/jpeg", // Default, could be enhanced to detect actual type
		FileSize:         fileInfo.Size(),
		UploadedByUserID: userID,
		CreatedAt:        fileInfo.ModTime(),
	}, nil
}

// UpdateUserAvatar handles the upload of a new avatar and replaces the old one
func (u *userAvatarUseCase) UpdateUserAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (*entity.Media, error) {
	log.Printf("find the user avatar usecase")
	if file.Size > maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of 5MB")
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("unsupported file type: %s", contentType)
	}

	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(avatarStoragePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Delete old avatar if exists
	oldURI, err := u.userAvatarRepo.ReadAvatarURI(ctx, userID)
	if err == nil && oldURI != "" {
		oldFilePath := strings.TrimPrefix(oldURI, "/")
		os.Remove(oldFilePath)
	}

	// Generate filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".jpg" // default extension
	}
	filename := fmt.Sprintf("%s%s", userID, ext)
	filePath := filepath.Join(avatarStoragePath, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Determine full URL for avatar
	baseURL := os.Getenv("UPLOAD_BASE_URL")
	// baseURL := "C:/Users/learn/Desktop/A2SV-Backend-Blog-Starter-Project"
	// uploadBaseURL := os.Getenv("UPLOAD_BASE_URL")
	avatarPath := fmt.Sprintf("/uploads/avatar_uri/%s", filename)

	// Ensure baseURL does not already include a trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Ensure avatarPath starts with a single slash
	if !strings.HasPrefix(avatarPath, "/") {
		avatarPath = "/" + avatarPath
	}

	avatarURI := fmt.Sprintf("%s%s", baseURL, avatarPath)
	// Update avatar URI in repository
	log.Printf("Updating avatar URI for user %s: %s", userID, avatarURI)
	if err := u.userAvatarRepo.UpdateAvatarURI(ctx, userID, avatarURI); err != nil {
		// Clean up file if database update fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to update avatar URI: %w", err)
	}

	// Create media entity
	media := &entity.Media{
		ID:               filename,
		FileName:         file.Filename,
		URL:              avatarURI,
		MimeType:         contentType,
		FileSize:         file.Size,
		UploadedByUserID: userID,
		CreatedAt:        time.Now(),
	}

	return media, nil
}

// DeleteUserAvatar removes the link to the avatar from the user's record
func (u *userAvatarUseCase) DeleteUserAvatar(ctx context.Context, userID string) error {
	// Get current avatar URI
	uri, err := u.userAvatarRepo.ReadAvatarURI(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to read avatar URI: %w", err)
	}

	if uri != "" {
		// Delete the actual file
		filePath := strings.TrimPrefix(uri, "/")
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}

	// Remove URI from repository
	return u.userAvatarRepo.DeleteAvatarURI(ctx, userID)
}
