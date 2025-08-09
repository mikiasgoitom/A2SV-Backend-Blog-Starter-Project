package usecase

import (
	"context"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// UserRepository defines the interface for user data persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUserPassword(ctx context.Context, id string, hashedPassword string) error
}

// BlogRepository interface defines the methods for blog data persistence.
type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *entity.Blog) error
	GetBlogByID(ctx context.Context, blogID string) (*entity.Blog, error)
	GetBlogs(ctx context.Context, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error)
	UpdateBlog(ctx context.Context, blog *entity.Blog) error
	DeleteBlog(ctx context.Context, blogID string) error
	SearchBlogs(ctx context.Context, query string, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error)
	GetBlogsByTags(ctx context.Context, tagIDs []string, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error)
	GetTrendingBlogs(ctx context.Context, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error)
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
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Validator defines the interface for generic input validation.
type Validator interface {
	ValidateEmail(email string) error
	ValidatePasswordStrength(password string) error
}
