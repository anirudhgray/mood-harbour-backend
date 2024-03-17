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

// GetUserMoodEntries handles getting mood entries for a user.
func (mc *MoodController) GetUserMoodEntries(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*models.User).ID
	moodTypeStr := c.Query("mood_type")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var moodType *models.MoodType
	if moodTypeStr != "" {
		moodTypeInt, err := strconv.Atoi(moodTypeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid-mood-type", "message": "Invalid mood type."})
			return
		}
		moodType = new(models.MoodType)
		*moodType = models.MoodType(moodTypeInt)
	}

	var startDate, endDate *string
	if startDateStr != "" {
		startDate = &startDateStr
	}
	if endDateStr != "" {
		endDate = &endDateStr
	}

	moodEntries, err := mc.moodService.GetUserMoodEntries(userID, moodType, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, moodEntries)
}

// GetSingleUserMoodEntry handles getting a single mood entry for a user.
func (mc *MoodController) GetSingleUserMoodEntry(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*models.User).ID
	moodIDStr := c.Param("id")
	moodID, err := strconv.ParseUint(moodIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid-id", "message": "Invalid mood ID."})
		return
	}

	moodEntry, err := mc.moodService.GetSingleUserMoodEntry(userID, uint(moodID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get-error", "message": err.Error()})
		return
	}

	if moodEntry.Mood.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "message": "You are not the owner of this mood entry."})
		return
	}

	c.JSON(http.StatusOK, moodEntry)
}

// DeleteMoodEntry handles deleting a mood entry.
func (mc *MoodController) DeleteMoodEntry(c *gin.Context) {
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

	err = mc.moodService.DeleteMoodEntry(userID, uint(moodID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mood entry deleted successfully."})
}

// CreateGenericAttribute handles creating a generic attribute.
func (mc *MoodController) CreateGenericAttribute(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*models.User).ID
	var attributeData struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&attributeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Improper JSON."})
		return
	}

	err := mc.moodService.CreateNewAttribute(attributeData.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "creation-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attribute created successfully."})
}

// GetAttributes handles getting all attributes.
func (mc *MoodController) GetGenericAttributes(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*models.User).ID
	attributes, err := mc.moodService.GetGenericAttributes(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attributes)
}
