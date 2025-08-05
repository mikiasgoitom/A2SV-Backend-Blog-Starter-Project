package dto

// CreateUserRequest is the DTO for creating a new user.
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32,containsuppercase,containslowercase,containsdigit,containssymbol"`
	FirstName string `json:"firstname" binding:"required,min=3,max=50"`
	LastName string `json:"lastname" binding:"required,min=3,max=50"`
}

// LoginRequest is the DTO for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest is the DTO for user registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

// UpdateUserRequest is the DTO for updating user profile.
type UpdateUserRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=32"`
	FirstName *string `json:"firstname,omitempty" binding:"omitempty,max=50"`
	LastName  *string `json:"lastname,omitempty" binding:"omitempty,max=50"`
	AvatarURL *string `json:"avatar_url,omitempty" binding:"omitempty,url"`
}

// ForgotPasswordRequest is the DTO for requesting password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest is the DTO for resetting password.
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

// VerifyEmailRequest is the DTO for verifying email.
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// ResendVerificationRequest is the DTO for resending verification email.
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RefreshTokenRequest is the DTO for refreshing access tokens.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
