package mocks

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/dto"
)

// MockUserUsecase is a mock implementation of the UserUsecase interface
type MockUserUsecase struct {
	// Control mock behavior
	ShouldFailCreateUser            bool
	ShouldFailVerifyEmail           bool
	ShouldFailResendVerificationEmail bool
	ShouldFailLogin                 bool
	ShouldFailGetByID               bool
	ShouldFailGetByEmail            bool
	ShouldFailUpdateUser            bool
	ShouldFailForgotPassword        bool
	ShouldFailResetPassword         bool
	ShouldFailRefreshToken          bool
	ShouldFailLogout                bool
	
	// Return values
	MockUser         entity.User
	MockAccessToken  string
	MockRefreshToken string
}

func NewMockUserUsecase() *MockUserUsecase {
	return &MockUserUsecase{
		MockUser: entity.User{
			ID:       uuid.New(),
			Username: "testuser",
			Email:    "test@example.com",
			Role:     entity.UserRoleUser,
		},
		MockAccessToken:  "mock_access_token",
		MockRefreshToken: "mock_refresh_token",
	}
}

func (m *MockUserUsecase) CreateUser(ctx context.Context, user entity.User, password string) (entity.User, error) {
	if m.ShouldFailCreateUser {
		return entity.User{}, errors.New("user creation failed")
	}
	return m.MockUser, nil
}

func (m *MockUserUsecase) VerifyEmail(ctx context.Context, token string) error {
	if m.ShouldFailVerifyEmail {
		return errors.New("email verification failed")
	}
	return nil
}

func (m *MockUserUsecase) ResendVerificationEmail(ctx context.Context, email string) error {
	if m.ShouldFailResendVerificationEmail {
		return errors.New("resend verification failed")
	}
	return nil
}

func (m *MockUserUsecase) Login(ctx context.Context, email, password string) (entity.User, string, string, error) {
	if m.ShouldFailLogin {
		return entity.User{}, "", "", errors.New("login failed")
	}
	return m.MockUser, m.MockAccessToken, m.MockRefreshToken, nil
}

func (m *MockUserUsecase) GetByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	if m.ShouldFailGetByID {
		return entity.User{}, errors.New("user not found")
	}
	return m.MockUser, nil
}

func (m *MockUserUsecase) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	if m.ShouldFailGetByEmail {
		return entity.User{}, errors.New("user not found")
	}
	return m.MockUser, nil
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (entity.User, error) {
	if m.ShouldFailUpdateUser {
		return entity.User{}, errors.New("update user failed")
	}
	return m.MockUser, nil
}

func (m *MockUserUsecase) ForgotPassword(ctx context.Context, email string) error {
	if m.ShouldFailForgotPassword {
		return errors.New("forgot password failed")
	}
	return nil
}

func (m *MockUserUsecase) ResetPassword(ctx context.Context, token, password string) error {
	if m.ShouldFailResetPassword {
		return errors.New("reset password failed")
	}
	return nil
}

func (m *MockUserUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if m.ShouldFailRefreshToken {
		return "", "", errors.New("refresh token failed")
	}
	return m.MockAccessToken, m.MockRefreshToken, nil
}

func (m *MockUserUsecase) Logout(ctx context.Context, userID uuid.UUID) error {
	if m.ShouldFailLogout {
		return errors.New("logout failed")
	}
	return nil
}
