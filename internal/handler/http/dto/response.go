package dto

import "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"

// UserResponse is the DTO for a user.
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// LoginResponse is the DTO for a successful login.
type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

func ToUserResponse(user entity.User) UserResponse {
	return UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}
}
