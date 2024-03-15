package repository

import (
	"testing"

	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/utils/test_utils"
)

func TestVerificationEntryRepository_CreateVerificationEntry(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.VerificationEntry{})

	// Create the VerificationEntry Repository with the test database
	ver := NewVerificationEntryRepository()
	ver.db = db

	// Create a test verification entry
	verificationEntry := models.VerificationEntry{
		Email: "test@example.com",
		OTP:   "123456", // Change this to your desired OTP
	}

	// Test CreateVerificationEntry function
	if err := ver.CreateVerificationEntry(verificationEntry); err != nil {
		t.Errorf("CreateVerificationEntry returned an error: %v", err)
	}

	// Test GetVerificationEntryByEmail function to retrieve the created verification entry
	retrievedEntry, err := ver.GetVerificationEntryByEmail("test@example.com")
	if err != nil {
		t.Errorf("GetVerificationEntryByEmail returned an error: %v", err)
	}
	if retrievedEntry.Email != verificationEntry.Email {
		t.Errorf("Expected entry email: %s, got: %s", verificationEntry.Email, retrievedEntry.Email)
	}
	if retrievedEntry.OTP != verificationEntry.OTP {
		t.Errorf("Expected entry OTP: %s, got: %s", verificationEntry.OTP, retrievedEntry.OTP)
	}
}

func TestVerificationEntryRepository_DeleteVerificationEntry(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.VerificationEntry{})

	// Create the VerificationEntry Repository with the test database
	ver := NewVerificationEntryRepository()
	ver.db = db

	// Create a test verification entry
	verificationEntry := models.VerificationEntry{
		Email: "test@example.com",
		OTP:   "123456", // Change this to your desired OTP
	}

	// Create the verification entry in the database
	if err := ver.CreateVerificationEntry(verificationEntry); err != nil {
		t.Fatalf("Failed to create a test verification entry: %v", err)
	}

	// Retrieve the verification entry by Email
	retrievedEntry, err := ver.GetVerificationEntryByEmail("test@example.com")
	if err != nil {
		t.Errorf("GetVerificationEntryByEmail returned an error: %v", err)
	}

	// Delete the verification entry by email
	if err := ver.DeleteVerificationEntry(retrievedEntry.Email); err != nil {
		t.Errorf("DeleteVerificationEntry returned an error: %v", err)
	}

	// Attempt to retrieve the verification entry to confirm deletion
	_, err = ver.GetVerificationEntryByEmail("test@example.com")
	if err == nil {
		t.Error("Verification entry was not deleted, GetVerificationEntryByEmail returned no error")
	} else if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found' error, got: %v", err)
	}
}
