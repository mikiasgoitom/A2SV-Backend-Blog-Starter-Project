package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	usecasecontract "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase/contract"
)

// UserAvatarHandler provides HTTP handlers for managing user avatars
type UserAvatarHandler struct {
	userAvatarUseCase usecasecontract.IUserAvatarUseCase
}

// NewUserAvatarHandler creates a new instance of UserAvatarHandler
func NewUserAvatarHandler(useCase usecasecontract.IUserAvatarUseCase) *UserAvatarHandler {
	return &UserAvatarHandler{
		userAvatarUseCase: useCase,
	}
}

// CreateUserAvatar handles POST request to create a new avatar
func (h *UserAvatarHandler) CreateUserAvatar(c *gin.Context) {
	userID := c.Param("id")
	log.Printf("DEBUG: userID from URL param: %s", userID)

	// Parse multipart form
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	// Create user avatar
	media, err := h.userAvatarUseCase.CreateUserAvatar(c.Request.Context(), userID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": media})
}

// ReadUserAvatar handles GET request to retrieve avatar data
func (h *UserAvatarHandler) ReadUserAvatar(c *gin.Context) {
	userID := c.Param("id")
	log.Printf("DEBUG: userID from URL param: %s", userID)

	media, err := h.userAvatarUseCase.ReadUserAvatarMetadata(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": media})
}

// UpdateUserAvatar handles PUT request to update a user's avatar
func (h *UserAvatarHandler) UpdateUserAvatar(c *gin.Context) {
	userID := c.Param("id")
	log.Printf("DEBUG: userID from URL param: %s", userID)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	media, err := h.userAvatarUseCase.UpdateUserAvatar(c.Request.Context(), userID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": media})
}

// DeleteUserAvatar handles DELETE request to remove a user's avatar
func (h *UserAvatarHandler) DeleteUserAvatar(c *gin.Context) {
	userID := c.Param("id")
	log.Printf("DEBUG: userID from URL param: %s", userID)

	if err := h.userAvatarUseCase.DeleteUserAvatar(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar deleted successfully"})
}

// RegisterRoutes registers user avatar routes with the Gin router
func (h *UserAvatarHandler) RegisterRoutes(router *gin.RouterGroup) {
	avatarRoutes := router.Group("/:id/avatar")
	{
		avatarRoutes.POST("", h.CreateUserAvatar)
		avatarRoutes.GET("", h.ReadUserAvatar)
		avatarRoutes.PUT("", h.UpdateUserAvatar)
		avatarRoutes.DELETE("", h.DeleteUserAvatar)
	}
}
