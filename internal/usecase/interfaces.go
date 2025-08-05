package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// UserRepository defines the interface for user data persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, id string, user *entity.User) error
	UpdateUserPassword(ctx context.Context, id string, hashedPassword string) error
}

// TokenRepository provides methods for managing tokens in the database.
type TokenRepository interface {
	CreateToken(ctx context.Context, token *entity.Token) error
	GetTokenByUserID(ctx context.Context, userID string) (*entity.Token, error)
	DeleteToken(ctx context.Context, tokenID string) error
	UpdateToken(ctx context.Context, tokenID string, tokenHash string, expiry time.Time) error
}

// EmailVerificationTokenRepository defines methods for interacting with email verification token data persistence.
type EmailVerificationTokenRepository interface {
	CreateEmailVerificationToken(ctx context.Context, token *entity.EmailVerificationToken) error
	GetEmailVerificationTokenByUserID(ctx context.Context, userID string) (*entity.EmailVerificationToken, error)
	GetEmailVerificationTokenByTokenHash(ctx context.Context, tokenHash string) (*entity.EmailVerificationToken, error)
	DeleteEmailVerificationToken(ctx context.Context, id string) error
	UpdateEmailVerificationTokenUsedStatus(ctx context.Context, id string, used bool) error
}

// Hasher provides methods for securely hashing and verifying passwords and other strings.
type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	HashString(s string) string
	CheckHash(s, hash string) bool
}

// JWTService defines the interface for JWT operations.
type JWTService interface {
	GenerateAccessToken(userID string, role entity.UserRole) (string, error)
	GenerateRefreshToken(userID string, role entity.UserRole) (string, error)
	ParseAccessToken(token string) (*entity.Claims, error)
	ParseRefreshToken(token string) (*entity.Claims, error)
	GeneratePasswordResetToken(userID string) (string, error)
	ParsePasswordResetToken(token string) (*entity.Claims, error)
	GenerateEmailVerificationToken(userID string) (string, error)
	ParseEmailVerificationToken(token string) (*entity.Claims, error)
}

// MailService defines the interface for sending emails.
type MailService interface {
	SendActivationEmail(toEmail, username, activationLink string) error
	SendPasswordResetEmail(toEmail, username, resetLink string) error
}

// ConfigProvider defines the interface for accessing application configuration.
type ConfigProvider interface {
	GetSendActivationEmail() bool
	GetAppBaseURL() string
	GetRefreshTokenExpiry() time.Duration
	GetPasswordResetTokenExpiry() time.Duration
	GetEmailVerificationTokenExpiry() time.Duration
}

// AppLogger defines the interface for logging messages.
type AppLogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(message string)
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// func (a AppLogger) Warn(s string) {
// 	panic("unimplemented")
// }

// Validator defines the interface for generic input validation.
type Validator interface {
	ValidateEmail(email string) error
	ValidatePasswordStrength(password string) error
}

// UserUseCase defines the interface for user-related operations.
type IUserUseCase interface {
	Register(ctx context.Context, username, email, password, firstName, lastName string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (*entity.User, string, string, error)
	Authenticate(ctx context.Context, accessToken string) (*entity.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, resetToken, newPassword string) error
	VerifyEmail(ctx context.Context, token string) error
	Logout(ctx context.Context, refreshToken string) error
	PromoteUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	DemoteUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
}

// type BlogUseCase interface {
// 	CreateBlog(ctx context.Context, blog entity.Blog) (*entity.Blog, error)
// 	GetBlogByID(ctx context.Context, blogID uuid.UUID) (*entity.Blog, error)
// 	UpdateBlog(ctx context.Context, blogID, authorID uuid.UUID, title *string, content *string, slug *string, status *entity.BlogStatus, publishedAt *time.Time, featuredImageID *uuid.UUID, isDeleted *bool) (*entity.Blog, error)
// 	TrackBlogPopularity(ctx context.Context, blogID, userID uuid.UUID, action BlogAction) (viewCount, likeCount, dislikeCount, commentCount int, err error)
// 	DeleteBlog(ctx context.Context, blogID, userID uuid.UUID, isAdmin bool) (bool, error)

// 	GetBlogs(ctx context.Context, page, pageSize int, sortBy string, sortOrder SortOrder, dateFrom *time.Time, dateTo *time.Time) (blogs []entity.Blog, totalCount int, currentPage int, totalPages int, err error)
// 	SearchAndFilterBlogs(ctx context.Context, query string, page, pageSize int, searchBy string, tags []string, dateFrom *time.Time, dateTo *time.Time, minViews *int, minLikes *int, authorID *uuid.UUID) (blogs []entity.Blog, totalCount int, currentPage int, totalPages int, err error)

// }
