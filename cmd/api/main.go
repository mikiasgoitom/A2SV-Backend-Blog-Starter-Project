package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	handlerHttp "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/config"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/external_services"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/jwt"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/logger"
	passwordservice "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/password_service"
	randomgenerator "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/random_generator"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/repository/mongodb"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/uuidgen"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/validator"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get MongoDB URI and DB name from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGODB_DB_NAME environment variable not set")
	}

	// Establish MongoDB connection
	mongoClient, err := mongodb.NewMongoDBClient(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect()

	// Initialize email service
	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")
	smtpUsername := os.Getenv("EMAIL_USERNAME")
	smtpPassword := os.Getenv("EMAIL_APP_PASSWORD")
	smtpFrom := os.Getenv("EMAIL_FROM")

	// Register custom validators
	validator.RegisterCustomValidators()

	// Initialize Gin router
	router := gin.Default()

	// Dependency Injection: Repositories
	userCollection := mongoClient.Client.Database(dbName).Collection("users")
	userRepo := mongodb.NewMongoUserRepository(userCollection)
	tokenRepo := mongodb.NewTokenRepository(mongoClient.Client.Database(dbName).Collection("tokens"))
	blogRepo := mongodb.NewBlogRepository(mongoClient.Client.Database(dbName), userCollection)
	likeRepo := mongodb.NewLikeRepository(mongoClient.Client.Database(dbName))

	// Dependency Injection: Services
	hasher := passwordservice.NewHasher()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	jwtManager := jwt.NewJWTManager(jwtSecret)
	jwtService := jwt.NewJWTService(jwtManager)
	appLogger := logger.NewStdLogger()
	mailService := external_services.NewEmailService(smtpHost, smtpPort, smtpUsername, smtpPassword, smtpFrom)
	randomGenerator := randomgenerator.NewRandomGenerator()
	appValidator := validator.NewValidator()
	uuidGenerator := uuidgen.NewGenerator()
	appConfig := config.NewConfig()
	baseURL := appConfig.GetAppBaseURL()
	// Dependency Injection: Usecases
	emailUsecase := usecase.NewEmailVerificationUseCase(tokenRepo, userRepo, mailService, randomGenerator, uuidGenerator, baseURL)
	userUsecase := usecase.NewUserUsecase(userRepo, tokenRepo, emailUsecase, hasher, jwtService, mailService, appLogger, appConfig, appValidator, uuidGenerator, randomGenerator)

	blogUsecase := usecase.NewBlogUseCase(blogRepo, uuidGenerator, appLogger)

	// Create like usecase
	likeUsecase := usecase.NewLikeUsecase(likeRepo, blogRepo)

	// Setup API routes
	appRouter := handlerHttp.NewRouter(userUsecase, blogUsecase, likeUsecase, emailUsecase, userRepo, tokenRepo, hasher, jwtService, mailService, appLogger, appConfig, appValidator, uuidGenerator, randomGenerator)
	appRouter.SetupRoutes(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
