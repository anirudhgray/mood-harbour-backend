package controllers

import (
	"net/http"
	"strconv"

	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/utils/recommender"
	"github.com/gin-gonic/gin"
)

func GenerateRecommendations(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(*models.User)
	userID := user.ID

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	likes, dislikes, err := recommender.GetUserLikesDislikes(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	otherUsers, err := recommender.GetUsersWithSimilarInteractions(likes, dislikes, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	similarities, err := recommender.CalculateSimilaritiesWithOtherUsers(userID, otherUsers, likes, dislikes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unreviewed, err := recommender.GetUnreviewedResources(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recprobs := recommender.CalculateRecommendationProbabilities(userID, unreviewed, similarities)

	recommendedResources := recommender.GetRecommendedResourcesSortedAndPaginated(recprobs, page, 20)

	var message string
	if len(recommendedResources) > 0 {
		message = "Here are some recommendations for you."
	} else {
		message = "Could not get any recommendations for you. Maybe try going through the discover page and reviewing some resources?."
	}

	c.JSON(http.StatusOK, gin.H{"message": message, "recommendations": recommendedResources})
}
