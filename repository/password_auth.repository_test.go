package repository

import (
	"testing"

	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/GDGVIT/attendance-app-backend/utils/test_utils"
)

func TestPasswordAuthRepository_CreatePwdAuthItem(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.PasswordAuth{})

	// Create the PasswordAuth Repository with the test database
	par := NewPasswordAuthRepository()
	par.db = db

	// Create a test PasswordAuth item
	passwordAuth := models.PasswordAuth{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Test CreatePwdAuthItem function
	if err := par.CreatePwdAuthItem(&passwordAuth); err != nil {
		t.Errorf("CreatePwdAuthItem returned an error: %v", err)
	}

	// Test GetPwdAuthItemByEmail function to retrieve the created item
	retrievedPasswordAuth, err := par.GetPwdAuthItemByEmail("test@example.com")
	if err != nil {
		t.Errorf("GetPwdAuthItemByEmail returned an error: %v", err)
	}
	if retrievedPasswordAuth.Email != passwordAuth.Email {
		t.Errorf("Expected email: %s, got: %s", passwordAuth.Email, retrievedPasswordAuth.Email)
	}
}

func TestPasswordAuthRepository_DeletePwdAuthItemByEmail(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.PasswordAuth{})

	// Create the PasswordAuth Repository with the test database
	par := NewPasswordAuthRepository()
	par.db = db

	// Create a test PasswordAuth item
	passwordAuth := models.PasswordAuth{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Create the item in the database
	if err := par.CreatePwdAuthItem(&passwordAuth); err != nil {
		t.Fatalf("Failed to create a test PasswordAuth item: %v", err)
	}

	// Test DeletePwdAuthItemByEmail function
	if err := par.DeletePwdAuthItemByEmail("test@example.com"); err != nil {
		t.Errorf("DeletePwdAuthItemByEmail returned an error: %v", err)
	}

	// Attempt to retrieve the item to confirm deletion
	_, err = par.GetPwdAuthItemByEmail("test@example.com")
	if err == nil {
		t.Error("PasswordAuth item was not deleted, GetPwdAuthItemByEmail returned no error")
	} else if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found' error, got: %v", err)
	}
}

func TestPasswordAuthRepository_UpdatePwdAuthItem(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.PasswordAuth{})

	// Create the PasswordAuth Repository with the test database
	par := NewPasswordAuthRepository()
	par.db = db

	// Create a test PasswordAuth item
	passwordAuth := models.PasswordAuth{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Create the item in the database
	if err := par.CreatePwdAuthItem(&passwordAuth); err != nil {
		t.Fatalf("Failed to create a test PasswordAuth item: %v", err)
	}

	// Retrieve the item by email
	retrievedPasswordAuth, err := par.GetPwdAuthItemByEmail("test@example.com")
	if err != nil {
		t.Errorf("GetPwdAuthItemByEmail returned an error: %v", err)
	}

	// Update the PasswordAuth item
	retrievedPasswordAuth.Password = "updatedpassword"
	if err := par.UpdatePwdAuthItem(retrievedPasswordAuth); err != nil {
		t.Errorf("UpdatePwdAuthItem returned an error: %v", err)
	}

	// Retrieve the item again and check if the password is updated
	retrievedPasswordAuth, err = par.GetPwdAuthItemByEmail("test@example.com")
	if err != nil {
		t.Errorf("GetPwdAuthItemByEmail returned an error: %v", err)
	}
	if retrievedPasswordAuth.Password != "updatedpassword" {
		t.Errorf("Expected password: updatedpassword, got: %s", retrievedPasswordAuth.Password)
	}
}
