package controllers

import (
	"net/http"
	"strconv"

	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/services"
	"github.com/gin-gonic/gin"
)

type ResourceController struct {
	resourceService services.ResourceServiceInterface
}

func NewResourceController(resourceService services.ResourceServiceInterface) *ResourceController {
	return &ResourceController{resourceService: resourceService}
}

// CreateResourceEntry creates a new resource entry in the database.
func (mc *ResourceController) CreateResourceEntry(c *gin.Context) {
	var resourceData struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		URL       string `json:"url"`
		External  bool   `json:"external"`
		AdminPost bool   `json:"admin_post"`
	}

	if err := c.ShouldBindJSON(&resourceData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Improper JSON."})
		return
	}

	user, _ := c.Get("user")
	userID := user.(*models.User).ID

	resourceResponse, err := mc.resourceService.CreateResourceEntry(userID, resourceData.Title, resourceData.Content, resourceData.URL, resourceData.External, resourceData.AdminPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "creation-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resourceResponse)
}

// GetAllResources gets all the resources in the database.
func (mc *ResourceController) GetAllResources(c *gin.Context) {
	resourceResponses, err := mc.resourceService.GetAllResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resourceResponses)
}

// GetResourceByID gets a resource by its ID.
func (mc *ResourceController) GetResourceByID(c *gin.Context) {
	resourceIDStr := c.Param("id")
	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid-id", "message": "Invalid resource ID."})
		return
	}

	resourceResponse, err := mc.resourceService.GetResourceByID(uint(resourceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resourceResponse)
}

// DeleteResourceEntry handles deleting a resource entry. Only the owner or an admin can delete a resource.
func (mc *ResourceController) DeleteResourceEntry(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*models.User).ID
	resourceIDStr := c.Param("id")
	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid-id", "message": "Invalid resource ID."})
		return
	}

	err = mc.resourceService.DeleteResource(userID, uint(resourceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource entry deleted successfully."})
}

// AddReview adds a review to a resource.
func (mc *ResourceController) AddReview(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*models.User).ID
	resourceIDStr := c.Param("id")
	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid-id", "message": "Invalid resource ID."})
		return
	}

	var reviewData struct {
		Content string        `json:"content"`
		Rating  models.Rating `json:"rating"`
	}

	if err := c.ShouldBindJSON(&reviewData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Improper JSON."})
		return
	}

	err = mc.resourceService.AddReview(uint(resourceID), userID, reviewData.Content, reviewData.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "creation-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Review added successfully."})
}
