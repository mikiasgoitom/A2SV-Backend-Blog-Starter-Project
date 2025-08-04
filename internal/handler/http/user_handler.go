package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/dto"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// CreateUser handles user registration (signup)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := BindAndValidate(c, &req); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Please make sure to fill all fields with the correct format.")
		return
	}

	_, err := h.userUsecase.Register(c.Request.Context(), req.Username, req.Email,req.Password, req.FirstName, req.LastName)
	if err != nil {
		ErrorHandler(c, http.StatusConflict, err.Error())
		return
	}

	MessageHandler(c, http.StatusCreated, "User created successfully. Please check your email to verify your account.")
}

// VerifyEmail handles email verification
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := BindAndValidate(c, &req); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Bad Request, Please try again.")
		return
	}

	err := h.userUsecase.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid or expired verification token")
		return
	}

	MessageHandler(c, http.StatusOK, "Email verified successfully. You can now log in.")
}

// Login handles user authentication
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := BindAndValidate(c, &req); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Bad Request credentials or unverified email")
		return
	}

	user, accessToken, refreshToken, err := h.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		ErrorHandler(c, http.StatusUnauthorized, "Invalid credentials or unverified email")
		return
	}

	response := dto.LoginResponse{
		User:         dto.ToUserResponse(*user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	SuccessHandler(c, http.StatusOK, response)
}

// GetUser handles retrieving user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	user, err := h.userUsecase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		ErrorHandler(c, http.StatusNotFound, "User not found")
		return
	}

	SuccessHandler(c, http.StatusOK, dto.ToUserResponse(*user))
}

// GetCurrentUser handles retrieving the current authenticated user
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.userUsecase.GetUserByID(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		ErrorHandler(c, http.StatusNotFound, "User not found")
		return
	}

	SuccessHandler(c, http.StatusOK, dto.ToUserResponse(*user))
}

// UpdateUser handles updating user profile
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.UpdateUserRequest
	if err := BindAndValidate(c, &req); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid or Bad request")

		return
	}

	updates := updateUserRequestToMap(req)
	updatedUser, err := h.userUsecase.UpdateProfile(c.Request.Context(), userID.(uuid.UUID), updates)
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	SuccessHandler(c, http.StatusOK, dto.ToUserResponse(*updatedUser))
}

// ForgotPassword handles password reset request
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := BindAndValidate(c, &req); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid or Bad request")

		return
	}

	err := h.userUsecase.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		MessageHandler(c, http.StatusOK, "If an account with that email exists, a password reset link has been sent")
		return
	}

	MessageHandler(c, http.StatusOK, "If an account with that email exists, a password reset link has been sent")
}

// ResetPassword handles password reset with token
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := BindAndValidate(c, &req); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid or Bad request")
		return
	}

	err := h.userUsecase.ResetPassword(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}

	MessageHandler(c, http.StatusOK, "Password reset successfully")
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		ErrorHandler(c, http.StatusBadRequest, "Refresh token required")
		return
	}

	// Remove "Bearer " prefix if present
	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	newAccessToken, newRefreshToken, err := h.userUsecase.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		ErrorHandler(c, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	response := gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	SuccessHandler(c, http.StatusOK, response)
}

// Logout handles user logout
func (h *UserHandler) Logout(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")

	err := h.userUsecase.Logout(c.Request.Context(), refreshToken)
	if err != nil {
		ErrorHandler(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	MessageHandler(c, http.StatusOK, "Logged out successfully")
}

func updateUserRequestToMap(req dto.UpdateUserRequest) (map[string] interface{}){
	updates := make(map[string]interface{})

	updates["username"] = req.Username
	updates["firstName"] = req.FirstName
	updates["lastName"] = req.LastName

	return updates
}