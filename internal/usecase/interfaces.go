package usecase

import (
	"context"

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

// UserUseCase defines the interface for user-related operations.
type IUserUseCase interface {
	Register(ctx context.Context, username, email, password, firstName, lastName string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (*entity.User, string, string, error)
	Authenticate(ctx context.Context, accessToken string) (*entity.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, verifier, resetToken, newPassword string) error
	VerifyEmail(ctx context.Context, token string) error
	Logout(ctx context.Context, refreshToken string) error
	PromoteUser(ctx context.Context, userID string) (*entity.User, error)
	DemoteUser(ctx context.Context, userID string) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) (*entity.User, error)
	GetUserByID(ctx context.Context, userID string) (*entity.User, error)
	LoginWithOAuth(ctx context.Context, firstName, lastName, email string) (string, string, error)
}
