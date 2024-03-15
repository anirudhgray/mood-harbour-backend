package middleware

import (
	"net/http"

	"github.com/GDGVIT/attendance-app-backend/infra/logger"
	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/GDGVIT/attendance-app-backend/repository"
	"github.com/GDGVIT/attendance-app-backend/utils/token"
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
