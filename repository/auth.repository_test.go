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
