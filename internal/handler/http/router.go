package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/middleware"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase"
	usecasecontract "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase/contract"
)

type Router struct {
	userHandler       *UserHandler
	userUsecase       *usecase.UserUsecase
	userAvatarUsecase usecasecontract.IUserAvatarUseCase
	jwtService        usecase.JWTService
}

func NewRouter(
	userUsecase *usecase.UserUsecase,
	avatarUsecase usecasecontract.IUserAvatarUseCase,
	jwtService usecase.JWTService,
) *Router {
	return &Router{
		userHandler:       NewUserHandler(userUsecase),
		userUsecase:       userUsecase,
		userAvatarUsecase: avatarUsecase,
		jwtService:        jwtService,
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	// Add a global logger middleware to see all incoming requests
	router.Use(gin.Logger())
	log.Println("DEBUG: Router setup initiated...")

	// API v1 routes
	v1 := router.Group("/api/v1")
	log.Println("DEBUG: Route group /api/v1 created")

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

	// Public user routes and avatar subroutes
	users := v1.Group("/users")
	{
		users.GET("/:id", r.userHandler.GetUser)
		// Avatar subroutes under /users/:id/avatar
		avatarHandler := NewUserAvatarHandler(r.userAvatarUsecase)
		log.Println("DEBUG: Registering avatar routes...")
		log.Println("DEBUG: Registering route /api/v1/users/:id/avatar")
		avatarHandler.RegisterRoutes(users)
	}

	// Protected routes (authentication required)
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleWare(r.jwtService, r.userUsecase))
	{
		log.Println("DEBUG: Setting up protected routes...")
		// Current user routes
		me := protected.Group("/me")
		{
			me.GET("", r.userHandler.GetCurrentUser)
			me.PUT("", r.userHandler.UpdateUser)
		}
	}

	// Logout route (no authentication required just accept the refresh token from the request body and invalidate the user session)
	v1.POST("/logout", r.userHandler.Logout)
}
