package middleware

import (
	"net/http"

	"github.com/anirudhgray/mood-harbour-backend/infra/logger"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/repository"
	"github.com/anirudhgray/mood-harbour-backend/utils/token"
	"github.com/gin-gonic/gin"
)

// BaseAuthMiddleware checks if the user is authenticated
func BaseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := token.ValidateToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "auth", "message": "Please login to continue."})
			// logger.Errorf("Auth Middleware Error: %v", err)
			c.Abort()
			return
		}

		var user models.User

		userRepo := repository.NewUserRepository()
		user, err = userRepo.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user-fetch-error", "message": "Internal error while fetching user."})
			logger.Errorf("Fetching Authenticated User Error: %v", err)
			c.Abort()
			return
		}

		c.Set("user", &user)

		c.Next()
	}
}

// AdminAuthMiddleware checks if the user is an admin
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		userData := user.(*models.User)

		if !userData.Admin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "auth", "message": "You are not authorized to access this resource."})
			c.Abort()
			return
		}

		c.Next()
	}
}
