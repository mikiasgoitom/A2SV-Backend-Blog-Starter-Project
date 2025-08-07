package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase"
)

type InteractionHandler struct {
	interactionUsecase usecase.IInteractionUseCase
}

func NewInteractionHandler(interactionUsecase usecase.IInteractionUseCase) *InteractionHandler {
	return &InteractionHandler{
		interactionUsecase: interactionUsecase,
	}
}

func (h *InteractionHandler) LikeBlogHandler(c *gin.Context) {
	blogID := c.Param("blogID")
	userID, exists := c.Get("userID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	fmt.Print("handler")

	userIDStr, ok := userID.(string)
	if !ok {
		ErrorHandler(c, http.StatusBadRequest, "Invalid user ID format in token")
		return
	}

	err := h.interactionUsecase.LikeBlog(c.Request.Context(), blogID, userIDStr)
	if err != nil {
		ErrorHandler(c, http.StatusInternalServerError, "Failed to like blog")
		return
	}

	SuccessHandler(c, http.StatusOK, "Blog liked successfully")
}

func (h *InteractionHandler) UnlikeBlogHandler(c *gin.Context) {
	blogID := c.Param("blogID")
	userID, exists := c.Get("userID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ErrorHandler(c, http.StatusBadRequest, "Invalid user ID format in token")
		return
	}

	err := h.interactionUsecase.UnlikeBlog(c.Request.Context(), blogID, userIDStr)
	if err != nil {
		ErrorHandler(c, http.StatusInternalServerError, "Failed to unlike blog")
		return
	}

	SuccessHandler(c, http.StatusOK, "Blog unliked successfully")
}
