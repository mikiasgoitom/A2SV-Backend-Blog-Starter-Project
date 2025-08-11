package middleware

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ValidateCommentContent validates comment content
func ValidateCommentContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData struct {
			Content string `json:"content"`
		}

		// Store the original request body for later use
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Unable to read request body",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Reset the request body so the handler can read it again
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Reset the request body again for the handler
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Required and non-empty validation
		content := strings.TrimSpace(requestData.Content)
		if content == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Content is required",
				"details": "Comment content cannot be empty",
			})
			c.Abort()
			return
		}

		// Length validation (1-1000 characters)
		if utf8.RuneCountInString(content) < 1 || utf8.RuneCountInString(content) > 1000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid content length",
				"details": "Comment content must be between 1 and 1000 characters",
			})
			c.Abort()
			return
		}

		// Basic XSS protection - remove HTML tags
		htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
		sanitizedContent := htmlTagRegex.ReplaceAllString(content, "")

		// Basic profanity filtering (simple implementation)
		profanityWords := []string{"spam", "abuse", "inappropriate"} // Add more as needed
		lowerContent := strings.ToLower(sanitizedContent)
		for _, word := range profanityWords {
			if strings.Contains(lowerContent, word) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Content contains inappropriate language",
					"details": "Please review your comment content",
				})
				c.Abort()
				return
			}
		}

		// Set sanitized content and validation flags
		c.Set("sanitized_content", sanitizedContent)
		c.Set("original_content", content)
		c.Set("content_validated", true)
		c.Next()
	}
}

// ValidateUUIDParam validates UUID parameters
func ValidateUUIDParam(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramValue := c.Param(paramName)
		if paramValue == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Missing parameter",
				"details": paramName + " is required",
			})
			c.Abort()
			return
		}

		if _, err := uuid.Parse(paramValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid UUID format",
				"details": paramName + " must be a valid UUID",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateCommentOwnership middleware to check if user owns the comment or is admin
func ValidateCommentOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware should be used after AuthMiddleware
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		userRole, roleExists := c.Get("userRole")
		if !roleExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found",
			})
			c.Abort()
			return
		}

		// If user is admin, allow access
		if userRole.(string) == "admin" {
			c.Next()
			return
		}

		// For regular users, we'll need to check ownership in the handler
		// This middleware just ensures the user is authenticated
		c.Set("user_id", userID.(string))
		c.Set("user_role", userRole.(string))
		c.Next()
	}
}

// CORSMiddleware handles CORS for frontend requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Configure as needed
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
