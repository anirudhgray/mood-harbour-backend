package test_utils

import (
	"github.com/anirudhgray/mood-harbour-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, error) {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// automigrate user and authprovider models
	db.AutoMigrate(&models.User{}, &models.AuthProvider{}, &models.DeletionConfirmation{}, &models.VerificationEntry{}, &models.ForgotPassword{}, &models.PasswordAuth{}, &models.Mood{}, &models.MoodAttribute{}, &models.Attribute{}, &models.Resource{}, &models.Review{})
	return db, nil
}
