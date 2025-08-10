package http

import (
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/middleware"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase"
	usecasecontract "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase/contract"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	userHandler        *UserHandler
	blogHandler        *BlogHandler
	emailHandler       *EmailHandler
	interactionHandler *InteractionHandler
	userUsecase        *usecase.UserUsecase
	jwtService         usecase.JWTService
	authHandler        *AuthHandler
}

func NewRouter(userUsecase usecasecontract.IUserUseCase, blogUsecase usecase.IBlogUseCase, likeUsecase *usecase.LikeUsecase, emailVerUC usecasecontract.IEmailVerificationUC, userRepo contract.IUserRepository, tokenRepo contract.ITokenRepository, hasher contract.IHasher, jwtService usecase.JWTService, mailService contract.IEmailService, logger usecasecontract.IAppLogger, config usecasecontract.IConfigProvider, validator usecasecontract.IValidator, uuidGen contract.IUUIDGenerator, randomGen contract.IRandomGenerator) *Router {
	baseURL := config.GetAppBaseURL()
	return &Router{
		userHandler:        NewUserHandler(userUsecase),
		blogHandler:        NewBlogHandler(blogUsecase),
		emailHandler:       NewEmailHandler(emailVerUC, userRepo),
		interactionHandler: NewInteractionHandler(likeUsecase),
		userUsecase:        usecase.NewUserUsecase(userRepo, tokenRepo, emailVerUC, hasher, jwtService, mailService, logger, config, validator, uuidGen, randomGen),
		jwtService:         jwtService,
		authHandler:        NewAuthHandler(userUsecase, baseURL),
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// rate limiter configuration
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	lmt.SetMessage("Too many requests, please try again later.")
	router.Use(middleware.RateLimiter(lmt))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/api/v1/metrics", gin.WrapH(promhttp.Handler()))
	// API v1 routes
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.userHandler.CreateUser)
		auth.POST("/login", r.userHandler.Login)
		auth.GET("/verify-email", r.emailHandler.HandleVerifyEmailToken)
		auth.POST("/forgot-password", r.userHandler.ForgotPassword)
		auth.POST("/reset-password", r.userHandler.ResetPassword)
		auth.POST("/refresh-token", r.userHandler.RefreshToken)

		auth.POST("/request-verification-email", r.emailHandler.HandleRequestEmailVerification)

		// Google OAuth endpoints
		auth.GET("/google/login", r.authHandler.HandleGoogleLogin)
		auth.GET("/google/callback", r.authHandler.HandleGoogleCallback)
	}

	// Public user routes
	users := v1.Group("/users")
	{
		users.GET("/:id", r.userHandler.GetUser)
	}

	// Public blog routes
	blogs := v1.Group("/blogs")
	{
		blogs.GET("", r.blogHandler.GetBlogsHandler)
		blogs.GET("/search", r.blogHandler.SearchAndFilterBlogsHandler)
		blogs.GET("/popular", r.blogHandler.GetPopularBlogsHandler)
		blogs.GET("/:slug", r.blogHandler.GetBlogDetailHandler)
	}

	// Protected routes (authentication required)
	protected := v1.Group("/")
	// Cast jwtService to the expected JWTManager type if needed
	protected.Use(middleware.AuthMiddleWare(r.jwtService, r.userUsecase))
	{
		// Current user routes
		protected.GET("/me", r.userHandler.GetCurrentUser)
		protected.PUT("/me", r.userHandler.UpdateUser)

		// Blog routes
		protected.POST("/blogs", r.blogHandler.CreateBlogHandler)
		protected.PUT("/blogs/:blogID", r.blogHandler.UpdateBlogHandler)
		protected.DELETE("/blogs/:blogID", r.blogHandler.DeleteBlogHandler)

		// Interaction routes
		protected.POST("/blogs/:blogID/like", r.interactionHandler.LikeBlogHandler)
		protected.DELETE("/blogs/:blogID/like", r.interactionHandler.UnlikeBlogHandler)
		protected.POST("/blogs/:blogID/dislike", r.interactionHandler.DislikeBlogHandler)
		protected.DELETE("/blogs/:blogID/dislike", r.interactionHandler.UndislikeBlogHandler)
		protected.POST("/blogs/:blogID/view", r.blogHandler.TrackBlogViewHandler)

	}

	// Logout route (no authentication required just accept the refresh token from the request body and invalidate the user session)
	v1.POST("/logout", r.userHandler.Logout)
}
