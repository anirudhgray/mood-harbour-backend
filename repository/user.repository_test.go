package repository

import (
	"testing"

	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/GDGVIT/attendance-app-backend/utils/test_utils"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.User{})

	// Create the User Repository with the test database
	ur := NewUserRepository()
	ur.db = db

	// Create a test user
	user := models.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Test CreateUser function
	if err := ur.CreateUser(user); err != nil {
		t.Errorf("CreateUser returned an error: %v", err)
	}

	// Test GetUserByEmail function to retrieve the created user
	retrievedUser, err := ur.GetUserByEmail("test@example.com")
	if err != nil {
		t.Errorf("GetUserByEmail returned an error: %v", err)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected user email: %s, got: %s", user.Email, retrievedUser.Email)
	}
}

func TestUserRepository_GetUserByEmail_NotFound(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.User{})

	// Create the User Repository with the test database
	ur := NewUserRepository()
	ur.db = db

	// Test GetUserByEmail function for a non-existent user
	_, err = ur.GetUserByEmail("test@example.com")
	if err == nil {
		t.Error("GetUserByEmail expected to return an error for non-existent user, but it didn't")
	}
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.User{})

	// Create the User Repository with the test database
	ur := NewUserRepository()
	ur.db = db

	// Create a test user
	user := models.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create the user in the database
	if err := ur.CreateUser(user); err != nil {
		t.Fatalf("Failed to create a test user: %v", err)
	}

	// Retrieve the user by Email
	retrievedUser, err := ur.GetUserByEmail(user.Email)
	if err != nil {
		t.Errorf("GetUserByEmail returned an error: %v", err)
	}

	// Delete the user by ID
	if err := ur.DeleteUserByID(retrievedUser.ID); err != nil {
		t.Errorf("DeleteUserByID returned an error: %v", err)
	}

	// Attempt to retrieve the user to confirm deletion
	_, err = ur.GetUserByID(user.ID)
	if err == nil {
		t.Error("User was not deleted, GetUserByID returned no error")
	} else if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found' error, got: %v", err)
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.User{})

	// Create the User Repository with the test database
	ur := NewUserRepository()
	ur.db = db

	// Create a test user
	user := models.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create the user in the database
	if err := ur.CreateUser(user); err != nil {
		t.Fatalf("Failed to create a test user: %v", err)
	}

	// Retrieve the user by Email
	retrievedUserByEmail, err := ur.GetUserByEmail(user.Email)
	if err != nil {
		t.Errorf("GetUserByEmail returned an error: %v", err)
	}

	// Retrieve the user by ID
	retrievedUser, err := ur.GetUserByID(retrievedUserByEmail.ID)
	if err != nil {
		t.Errorf("GetUserByID returned an error: %v", err)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected user email: %s, got: %s", user.Email, retrievedUser.Email)
	}
}

func TestUserRepository_VerifyUserEmail(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.User{}) // Drop the table after testing

	// Create the User Repository with the test database
	ur := NewUserRepository()
	ur.db = db

	// Create a test user
	user := models.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create the user in the database
	if err := ur.CreateUser(user); err != nil {
		t.Fatalf("Failed to create a test user: %v", err)
	}

	// Verify the user's email
	if err := ur.VerifyUserEmail(user.Email); err != nil {
		t.Errorf("VerifyUserEmail returned an error: %v", err)
	}

	// Retrieve the user and check if their email is verified
	retrievedUser, err := ur.GetUserByEmail(user.Email)
	if err != nil {
		t.Errorf("GetUserByEmail returned an error: %v", err)
	}
	if !retrievedUser.Verified {
		t.Error("Email verification failed, user's email is not marked as verified")
	}
}

func TestUserRepository_SaveUser(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.User{}) // Drop the table after testing

	// Create the User Repository with the test database
	ur := NewUserRepository()
	ur.db = db

	// Create a test user
	user := models.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create the user in the database
	if err := ur.CreateUser(user); err != nil {
		t.Fatalf("Failed to create a test user: %v", err)
	}

	// Retrieve the user by Email
	retrievedUser, err := ur.GetUserByEmail(user.Email)
	if err != nil {
		t.Errorf("GetUserByEmail returned an error: %v", err)
	}

	// Update the user's name
	retrievedUser.Name = "Updated User"
	if err := ur.SaveUser(retrievedUser); err != nil {
		t.Errorf("SaveUser returned an error: %v", err)
	}

	// Retrieve the user and check if their name is updated
	retrievedUser, err = ur.GetUserByEmail(user.Email)
	if err != nil {
		t.Errorf("GetUserByEmail returned an error: %v", err)
	}
	if retrievedUser.Name != "Updated User" {
		t.Errorf("Expected user name: Updated User, got: %s", retrievedUser.Name)
	}
}
