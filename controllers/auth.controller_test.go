package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anirudhgray/mood-harbour-backend/mocks"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserController_RegisterUser(t *testing.T) {
	// Initialize the mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock AuthService
	mockService := mocks.NewMockAuthServiceInterface(ctrl)

	// Create a new Gin router
	r := gin.Default()

	// Initialize the user controller with the mock service
	userController := NewUserController(mockService)

	// Register routes for testing
	r.POST("/auth/register", userController.RegisterUser)

	// Create a test user input
	user := models.User{
		Email:        "test@test.com",
		Name:         "Test User",
		ProfileImage: "test.jpg",
		Verified:     false,
		Model: gorm.Model{
			ID:        1,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
	}

	// Mock the service's RegisterUser function
	mockService.EXPECT().RegisterUser(user.Email, user.Name, user.ProfileImage, "pwd").Return(user, nil)

	var registerData struct {
		Email        string `json:"email"`
		Name         string `json:"name"`
		Password     string `json:"password"`
		ProfileImage string `json:"profile_image"`
	}

	registerData.Email = user.Email
	registerData.Name = user.Name
	registerData.Password = "pwd"
	registerData.ProfileImage = user.ProfileImage

	// Create a test request
	reqBody, err := json.Marshal(registerData)
	assert.NoError(t, err)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)
	// You can also parse the response body to validate the result
	var responseUser struct {
		Email        string `json:"email"`
		Name         string `json:"name"`
		Password     string `json:"password"`
		ProfileImage string `json:"profile_image"`
	}
	err = json.NewDecoder(w.Body).Decode(&responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, responseUser.Email)
}
