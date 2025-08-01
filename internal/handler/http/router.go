package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase"
)

type Router struct {
	userHandler *UserHandler
}

func NewRouter(userUsecase usecase.UserUsecase) *Router {
	return &Router{
		userHandler: NewUserHandler(userUsecase),
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	// API v1 routes
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.userHandler.CreateUser)
		auth.POST("/login", r.userHandler.Login)
		auth.POST("/verify-email", r.userHandler.VerifyEmail)
		auth.POST("/resend-verification", r.userHandler.ResendVerification)
		auth.POST("/forgot-password", r.userHandler.ForgotPassword)
		auth.POST("/reset-password", r.userHandler.ResetPassword)
		auth.POST("/refresh-token", r.userHandler.RefreshToken)
	}

	// Public user routes
	users := v1.Group("/users")
	{
		users.GET("/:id", r.userHandler.GetUser)
	}

	// Protected routes (authentication required)
	// Note: Add authentication middleware when it's properly configured
	protected := v1.Group("/")
	// protected.Use(middleware.AuthMiddleWare(jwtManager, userUsecase))
	{
		// Current user routes
		protected.GET("/me", r.userHandler.GetCurrentUser)
		protected.PUT("/me", r.userHandler.UpdateUser)
		protected.POST("/logout", r.userHandler.Logout)
	}
}
