package routers

import (
	"net/http"

	"github.com/anirudhgray/mood-harbour-backend/controllers"
	"github.com/anirudhgray/mood-harbour-backend/repository"
	"github.com/anirudhgray/mood-harbour-backend/routers/middleware"
	"github.com/anirudhgray/mood-harbour-backend/services"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

	authProviderRepo := repository.NewAuthProviderRepository()
	forgotPasswordRepo := repository.NewForgotPasswordRepository()
	deletionConfirmationRepo := repository.NewDeletionConfirmationRepository()
	verificationRepo := repository.NewVerificationEntryRepository()
	passwordAuthRepo := repository.NewPasswordAuthRepository()
	userRepo := repository.NewUserRepository()
	moodRepo := repository.NewMoodRepository()

	emailService := services.NewEmailService(userRepo)
	authService := services.NewAuthService(
		authProviderRepo,
		verificationRepo,
		forgotPasswordRepo,
		deletionConfirmationRepo,
		passwordAuthRepo,
		userRepo,
		emailService,
	)
	moodService := services.NewMoodService(moodRepo)

	v1 := route.Group("/v1")

	auth := v1.Group("/auth") // Create an /auth/ group
	{
		userController := controllers.NewUserController(authService) // Create an instance of the UserController

		// Define the user registration route
		auth.POST("/register", userController.RegisterUser)

		// Define the user login route
		auth.POST("/login", userController.Login)

		// Verify user account by providing otp
		auth.POST("/verify", userController.VerifyEmail)

		// Request another verification email
		auth.GET("/request-verification", userController.RequestVerificationAgain)

		// Send forgot password request
		auth.GET("/forgot-password", userController.ForgotPasswordRequest)

		// Set forgotten password
		auth.POST("/set-forgotten-password", userController.SetNewPassword)

		// Test baseauth middleware
		auth.GET("/test-auth", middleware.BaseAuthMiddleware(), userController.TestAuth)

		// Reset password by logged in user
		auth.POST("/reset-password", middleware.BaseAuthMiddleware(), userController.ResetPassword)

		// Send account deletion request
		auth.GET("/request-delete-account", middleware.BaseAuthMiddleware(), userController.RequestDeletion)

		// Delete account
		auth.DELETE("/delete-account", middleware.BaseAuthMiddleware(), userController.DeleteAccount)

		// Google login
		auth.GET("/google/login", userController.GoogleLogin)

		// Google Callback
		auth.GET("/google/callback", userController.GoogleCallback)
	}

	mood := v1.Group("/mood", middleware.BaseAuthMiddleware())
	{
		moodController := controllers.NewMoodController(moodService)

		// Create a new mood entry
		mood.POST("/create", moodController.CreateMoodEntry)

		// Update a mood entry
		mood.PUT("/update/:id", moodController.UpdateUserMoodEntry)
	}

}
