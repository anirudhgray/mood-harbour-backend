package repository

import (
	"testing"

	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/GDGVIT/attendance-app-backend/utils/test_utils"
)

func TestDeletionConfirmationRepository_CreateDeletionConfirmation(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.DeletionConfirmation{})

	// Create the DeletionConfirmation Repository with the test database
	dcr := NewDeletionConfirmationRepository()
	dcr.db = db

	// Create a test deletion confirmation entry
	confirmation := models.DeletionConfirmation{
		Email: "test@example.com",
		OTP:   "123456",
	}

	// Test CreateDeletionConfirmation function
	if err := dcr.CreateDeletionConfirmation(confirmation); err != nil {
		t.Errorf("CreateDeletionConfirmation returned an error: %v", err)
	}

	// Test GetDeletionConfirmationByEmail function to retrieve the created entry
	retrievedConfirmation, err := dcr.GetDeletionConfirmationByEmail("test@example.com")
	if err != nil {
		t.Errorf("GetDeletionConfirmationByEmail returned an error: %v", err)
	}
	if retrievedConfirmation.Email != confirmation.Email {
		t.Errorf("Expected email: %s, got: %s", confirmation.Email, retrievedConfirmation.Email)
	}
}

func TestDeletionConfirmationRepository_DeleteDeletionConfirmationByEmail(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.DeletionConfirmation{})

	// Create the DeletionConfirmation Repository with the test database
	dcr := NewDeletionConfirmationRepository()
	dcr.db = db

	// Create a test deletion confirmation entry
	confirmation := models.DeletionConfirmation{
		Email: "test@example.com",
		OTP:   "123456",
	}

	// Create the entry in the database
	if err := dcr.CreateDeletionConfirmation(confirmation); err != nil {
		t.Fatalf("Failed to create a test deletion confirmation entry: %v", err)
	}

	// Test DeleteDeletionConfirmationByEmail function
	if err := dcr.DeleteDeletionConfirmationByEmail("test@example.com"); err != nil {
		t.Errorf("DeleteDeletionConfirmationByEmail returned an error: %v", err)
	}

	// Attempt to retrieve the entry to confirm deletion
	_, err = dcr.GetDeletionConfirmationByEmail("test@example.com")
	if err == nil {
		t.Error("Deletion confirmation entry was not deleted, GetDeletionConfirmationByEmail returned no error")
	} else if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found' error, got: %v", err)
	}
}
