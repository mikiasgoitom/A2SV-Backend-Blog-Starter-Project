package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/middleware"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase"
)

type Router struct {
	userHandler *UserHandler
	userUsecase *usecase.UserUsecase
	jwtService  usecase.JWTService
}

func NewRouter(userUsecase *usecase.UserUsecase, jwtService usecase.JWTService) *Router {
	return &Router{
		userHandler: NewUserHandler(userUsecase),
		userUsecase: userUsecase,
		jwtService:  jwtService,
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
	protected := v1.Group("/")
	// Cast jwtService to the expected JWTManager type if needed
	protected.Use(middleware.AuthMiddleWare(r.jwtService, r.userUsecase))
	{
		// Current user routes
		protected.GET("/me", r.userHandler.GetCurrentUser)
		protected.PUT("/me", r.userHandler.UpdateUser)
	}

	// Logout route (no authentication required just accept the refresh token from the request body and invalidate the user session)
	v1.POST("/logout", r.userHandler.Logout)
}
