package repository

import (
	"testing"

	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/utils/test_utils"
)

func TestForgotPasswordRepository_CreateForgotPassword(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.ForgotPassword{})

	// Create the ForgotPassword Repository with the test database
	fpr := NewForgotPasswordRepository()
	fpr.db = db

	// Create a test forgot password entry
	forgotPassword := models.ForgotPassword{
		Email: "test@example.com",
		OTP:   "123456",
	}

	// Test CreateForgotPassword function
	if err := fpr.CreateForgotPassword(forgotPassword); err != nil {
		t.Errorf("CreateForgotPassword returned an error: %v", err)
	}

	// Test GetForgotPasswordByEmail function to retrieve the created entry
	retrievedForgotPassword, err := fpr.GetForgotPasswordByEmail("test@example.com")
	if err != nil {
		t.Errorf("GetForgotPasswordByEmail returned an error: %v", err)
	}
	if retrievedForgotPassword.Email != forgotPassword.Email {
		t.Errorf("Expected email: %s, got: %s", forgotPassword.Email, retrievedForgotPassword.Email)
	}
}

func TestForgotPasswordRepository_DeleteForgotPasswordByEmail(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.ForgotPassword{})

	// Create the ForgotPassword Repository with the test database
	fpr := NewForgotPasswordRepository()
	fpr.db = db

	// Create a test forgot password entry
	forgotPassword := models.ForgotPassword{
		Email: "test@example.com",
		OTP:   "123456",
	}

	// Create the entry in the database
	if err := fpr.CreateForgotPassword(forgotPassword); err != nil {
		t.Fatalf("Failed to create a test forgot password entry: %v", err)
	}

	// Test DeleteForgotPasswordByEmail function
	if err := fpr.DeleteForgotPasswordByEmail("test@example.com"); err != nil {
		t.Errorf("DeleteForgotPasswordByEmail returned an error: %v", err)
	}

	// Attempt to retrieve the entry to confirm deletion
	_, err = fpr.GetForgotPasswordByEmail("test@example.com")
	if err == nil {
		t.Error("Forgot password entry was not deleted, GetForgotPasswordByEmail returned no error")
	} else if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found' error, got: %v", err)
	}
}
