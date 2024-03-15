package repository

import (
	"testing"

	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/GDGVIT/attendance-app-backend/utils/test_utils"
)

func TestAuthProviderRepository_GetAuthProviderByProviderKey(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.AuthProvider{}) // Drop the table after testing

	// Create the AuthProvider Repository with the test database
	ar := NewAuthProviderRepository()
	ar.db = db

	// Create a test AuthProvider
	authProvider := models.AuthProvider{
		ProviderKey:  "google_key",
		UserID:       1,
		ProviderName: "Google",
	}

	// Create the AuthProvider in the database
	if err := ar.CreateAuthProvider(authProvider); err != nil {
		t.Fatalf("Failed to create a test AuthProvider: %v", err)
	}

	// Retrieve the AuthProvider by its ProviderKey
	retrievedAuthProvider, err := ar.GetAuthProviderByProviderKey(authProvider.ProviderKey)
	if err != nil {
		t.Errorf("GetAuthProviderByProviderKey returned an error: %v", err)
	}
	if retrievedAuthProvider.ProviderKey != authProvider.ProviderKey {
		t.Errorf("Expected ProviderKey: %s, got: %s", authProvider.ProviderKey, retrievedAuthProvider.ProviderKey)
	}
}

func TestAuthProviderRepository_CreateAuthProvider(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.AuthProvider{}) // Drop the table after testing

	// Create the AuthProvider Repository with the test database
	ar := NewAuthProviderRepository()
	ar.db = db

	// Create a test AuthProvider
	authProvider := models.AuthProvider{
		ProviderKey:  "google_key",
		UserID:       1,
		ProviderName: "Google",
	}

	// Create the AuthProvider in the database
	if err := ar.CreateAuthProvider(authProvider); err != nil {
		t.Fatalf("Failed to create a test AuthProvider: %v", err)
	}

	retrievedAuthProviderByKey, err := ar.GetAuthProviderByProviderKey(authProvider.ProviderKey)
	if err != nil {
		t.Errorf("GetAuthProviderByProviderKey returned an error: %v", err)
	}

	// Retrieve the AuthProvider by its ID
	retrievedAuthProvider, err := ar.GetAuthProviderByID(retrievedAuthProviderByKey.ID)
	if err != nil {
		t.Errorf("GetAuthProviderByID returned an error: %v", err)
	}
	if retrievedAuthProvider.ProviderKey != authProvider.ProviderKey {
		t.Errorf("Expected ProviderKey: %s, got: %s", authProvider.ProviderKey, retrievedAuthProvider.ProviderKey)
	}
}

func TestAuthProviderRepository_UpdateAuthProvider(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.AuthProvider{}) // Drop the table after testing

	// Create the AuthProvider Repository with the test database
	ar := NewAuthProviderRepository()
	ar.db = db

	// Create a test AuthProvider
	authProvider := models.AuthProvider{
		ProviderKey:  "google_key",
		UserID:       1,
		ProviderName: "Google",
	}

	// Create the AuthProvider in the database
	if err := ar.CreateAuthProvider(authProvider); err != nil {
		t.Fatalf("Failed to create a test AuthProvider: %v", err)
	}

	retrievedAuthProviderByKey, err := ar.GetAuthProviderByProviderKey(authProvider.ProviderKey)
	if err != nil {
		t.Errorf("GetAuthProviderByProviderKey returned an error: %v", err)
	}

	// Update the AuthProvider's ProviderName
	retrievedAuthProviderByKey.ProviderName = "Updated Google"
	if err := ar.UpdateAuthProvider(retrievedAuthProviderByKey); err != nil {
		t.Errorf("UpdateAuthProvider returned an error: %v", err)
	}

	// Retrieve the AuthProvider and check if the ProviderName is updated
	retrievedAuthProvider, err := ar.GetAuthProviderByID(retrievedAuthProviderByKey.ID)
	if err != nil {
		t.Errorf("GetAuthProviderByID returned an error: %v", err)
	}
	if retrievedAuthProvider.ProviderName != "Updated Google" {
		t.Errorf("Expected ProviderName: Updated Google, got: %s", retrievedAuthProvider.ProviderName)
	}
}

func TestAuthProviderRepository_DeleteAuthProvider(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.AuthProvider{}) // Drop the table after testing

	// Create the AuthProvider Repository with the test database
	ar := NewAuthProviderRepository()
	ar.db = db

	// Create a test AuthProvider
	authProvider := models.AuthProvider{
		ProviderKey:  "google_key",
		UserID:       1,
		ProviderName: "Google",
	}

	// Create the AuthProvider in the database
	if err := ar.CreateAuthProvider(authProvider); err != nil {
		t.Fatalf("Failed to create a test AuthProvider: %v", err)
	}

	// Delete the AuthProvider by its ID
	if err := ar.DeleteAuthProvider(authProvider.ID); err != nil {
		t.Errorf("DeleteAuthProvider returned an error: %v", err)
	}

	// Attempt to retrieve the AuthProvider to confirm deletion
	_, err = ar.GetAuthProviderByID(authProvider.ID)
	if err == nil {
		t.Error("AuthProvider was not deleted, GetAuthProviderByID returned no error")
	} else if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found' error, got: %v", err)
	}
}

func TestAuthProviderRepository_DeleteAuthProviderByUserID(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.AuthProvider{}) // Drop the table after testing

	// Create the AuthProvider Repository with the test database
	ar := NewAuthProviderRepository()
	ar.db = db

	// Create a test user
	user := models.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create a test AuthProvider associated with the user
	authProvider := models.AuthProvider{
		ProviderKey:  "google_key",
		UserID:       user.ID,
		ProviderName: "Google",
	}

	// Create the AuthProvider in the database
	if err := ar.CreateAuthProvider(authProvider); err != nil {
		t.Fatalf("Failed to create a test AuthProvider: %v", err)
	}

	// Delete the AuthProvider by its UserID
	if err := ar.DeleteAuthProviderByUserID(user.ID); err != nil {
		t.Errorf("DeleteAuthProviderByUserID returned an error: %v", err)
	}

	// Attempt to retrieve the AuthProvider to confirm deletion
	_, err = ar.GetAuthProviderByID(authProvider.ID)
	if err == nil {
		t.Error("AuthProvider was not deleted, GetAuthProviderByID returned no error")
	} else if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found' error, got: %v", err)
	}
}
