package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/dto"
)

// UserUsecase defines the user use case interface
type UserUsecase interface {
	CreateUser(ctx context.Context, user entity.User, password string) (entity.User, error)
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, email string) error
	Login(ctx context.Context, email, password string) (entity.User, string, string, error)
	GetByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (entity.User, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, password string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}

// BlogUsecase defines the blog use case interface
type BlogUsecase interface {
	// Blog methods will be defined here
}

// AIUsecase defines the AI use case interface
type AIUsecase interface {
	// AI methods will be defined here
}
