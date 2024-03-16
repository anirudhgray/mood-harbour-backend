package controllers

import (
	"net/http"
	"strconv"

	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/services"
	"github.com/gin-gonic/gin"
)

type MoodController struct {
	moodService services.MoodServiceInterface
}

// NewMoodController creates a new MoodController
func NewMoodController(moodService services.MoodServiceInterface) *MoodController {
	return &MoodController{moodService: moodService}
}

// CreateMoodEntry handles mood entry creation.
func (mc *MoodController) CreateMoodEntry(c *gin.Context) {
	var moodData struct {
		MoodType   models.MoodType `json:"mood_type"`
		Notes      string          `json:"notes"`
		Attributes []string        `json:"attributes"`
	}

	if err := c.ShouldBindJSON(&moodData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Improper JSON."})
		return
	}

	user, _ := c.Get("user")
	userID := user.(*models.User).ID

	moodResponse, err := mc.moodService.CreateMoodEntry(moodData.MoodType, moodData.Notes, userID, moodData.Attributes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "creation-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, moodResponse)
}

// UpdateUserMoodEntry handles mood entry updates.
func (mc *MoodController) UpdateUserMoodEntry(c *gin.Context) {
	// check if current user is the owner of the mood entry
	user, _ := c.Get("user")
	userID := user.(*models.User).ID
	moodIDStr := c.Param("id")
	moodID, err := strconv.ParseUint(moodIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid-id", "message": "Invalid mood ID."})
		return
	}

	// get mood entry from database
	moodEntry, err := mc.moodService.GetSingleUserMoodEntry(userID, uint(moodID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get-error", "message": err.Error()})
		return
	}
	if moodEntry.Mood.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "message": "You are not the owner of this mood entry."})
		return
	}

	var moodData struct {
		MoodType   models.MoodType `json:"mood_type"`
		Notes      string          `json:"notes"`
		Attributes []string        `json:"attributes"`
	}

	if err := c.ShouldBindJSON(&moodData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Improper JSON."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid-id", "message": "Invalid mood ID."})
		return
	}

	err = mc.moodService.UpdateUserMoodEntry(uint(moodID), moodData.MoodType, moodData.Notes, moodData.Attributes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mood entry updated successfully."})
}
