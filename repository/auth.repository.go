package repository

import (
	"github.com/anirudhgray/mood-harbour-backend/infra/database" // Import your custom database package
	"github.com/anirudhgray/mood-harbour-backend/infra/logger"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"gorm.io/gorm"
)

// **UserRepository** //

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{database.DB}
}

type UserRepositoryInterface interface {
	CreateUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(userID uint) (models.User, error)
	VerifyUserEmail(email string) error
	SaveUser(user models.User) error
	DeleteUserByID(userID uint) error
}

func (ur *UserRepository) CreateUser(user models.User) error {
	if err := ur.db.Create(&user).Error; err != nil {
		logger.Errorf("DB: Error Creating User: %v", err)
		return err
	}
	return nil
}

func (ur *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorf("DB: User Record not found")
			return user, err // User not found
		}
		logger.Errorf("DB: Error Getting User By Email: %v", err)
		return user, err
	}
	return user, nil
}

// GetUserByID fetches a user by their ID
func (ur *UserRepository) GetUserByID(userID uint) (models.User, error) {
	var user models.User
	if err := ur.db.First(&user, userID).Error; err != nil {
		return user, err
	}
	return user, nil
}

// VerifyUserEmail verifies a user's email by updating the verification status
func (ur *UserRepository) VerifyUserEmail(email string) error {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	user.Verified = true
	if err := ur.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// SaveUser saves user model
func (ur *UserRepository) SaveUser(user models.User) error {
	if err := ur.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUserByID deletes a user by their ID
func (ur *UserRepository) DeleteUserByID(userID uint) error {
	var user models.User
	if err := ur.db.First(&user, userID).Error; err != nil {
		return err
	}

	if err := ur.db.Unscoped().Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

// for testing purposes
func (ur *UserRepository) SetDB(db *gorm.DB) {
	ur.db = db
}

// **VerificationEntryRepository** //

type VerificationEntryRepository struct {
	db *gorm.DB
}

func NewVerificationEntryRepository() *VerificationEntryRepository {
	return &VerificationEntryRepository{database.DB}
}

type VerificationRepositoryInterface interface {
	CreateVerificationEntry(verificationEntry models.VerificationEntry) error
	GetVerificationEntryByEmail(email string) (*models.VerificationEntry, error)
	DeleteVerificationEntry(email string) error
}

// CreateVerificationEntry creates a new verification entry
func (ver *VerificationEntryRepository) CreateVerificationEntry(verificationEntry models.VerificationEntry) error {
	if err := ver.db.Create(&verificationEntry).Error; err != nil {
		logger.Errorf("DB: Error Creating VerificationEntry: %v", err)
		return err
	}
	return nil
}

// GetVerificationEntryByEmail fetches a verification entry by email
func (ver *VerificationEntryRepository) GetVerificationEntryByEmail(email string) (*models.VerificationEntry, error) {
	var verificationEntry models.VerificationEntry
	if err := ver.db.Where("email = ?", email).First(&verificationEntry).Error; err != nil {
		logger.Errorf("DB: Error Getting VerificationEntry: %v", err)
		return nil, err
	}
	return &verificationEntry, nil
}

// DeleteVerificationEntry deletes a verification entry by email
func (ver *VerificationEntryRepository) DeleteVerificationEntry(email string) error {
	if err := ver.db.Where("email = ?", email).Unscoped().Delete(&models.VerificationEntry{}).Error; err != nil {
		logger.Errorf("DB: Error Deleting VerificationEntry: %v", err)
		return err
	}
	return nil
}

// **PasswordAuthRepository** //

type PasswordAuthRepository struct {
	db *gorm.DB
}

func NewPasswordAuthRepository() *PasswordAuthRepository {
	return &PasswordAuthRepository{database.DB}
}

type PasswordAuthRepositoryInterface interface {
	CreatePwdAuthItem(passwordAuth *models.PasswordAuth) error
	GetPwdAuthItemByEmail(email string) (models.PasswordAuth, error)
	UpdatePwdAuthItem(passwordAuth models.PasswordAuth) error
	DeletePwdAuthItem(id uint) error
	DeletePwdAuthItemByEmail(email string) error
}

// CreateUser creates a new PasswordAuth record.
func (par *PasswordAuthRepository) CreatePwdAuthItem(passwordAuth *models.PasswordAuth) error {
	return par.db.Create(passwordAuth).Error
}

// GetUserByEmail retrieves a PasswordAuth record by email.
func (par *PasswordAuthRepository) GetPwdAuthItemByEmail(email string) (models.PasswordAuth, error) {
	var user models.PasswordAuth
	if err := par.db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser updates an existing PasswordAuth record.
func (par *PasswordAuthRepository) UpdatePwdAuthItem(passwordAuth models.PasswordAuth) error {
	return par.db.Save(passwordAuth).Error
}

// DeleteUser deletes a PasswordAuth record by ID.
func (par *PasswordAuthRepository) DeletePwdAuthItem(id uint) error {
	return par.db.Unscoped().Delete(&models.PasswordAuth{}, id).Error
}

// DeletePwdAuthItemByEmail deletes a PasswordAuth record by email.
func (par *PasswordAuthRepository) DeletePwdAuthItemByEmail(email string) error {
	return par.db.Where("email = ?", email).Unscoped().Delete(&models.PasswordAuth{}).Error
}

// **ForgotPasswordRepository** //

type ForgotPasswordRepository struct {
	db *gorm.DB
}

func NewForgotPasswordRepository() *ForgotPasswordRepository {
	return &ForgotPasswordRepository{database.DB}
}

type ForgotPasswordRepositoryInterface interface {
	CreateForgotPassword(forgotPassword models.ForgotPassword) error
	GetForgotPasswordByEmail(email string) (*models.ForgotPassword, error)
	DeleteForgotPasswordByEmail(email string) error
}

// CreateForgotPassword creates a new forgot password entry
func (fpr *ForgotPasswordRepository) CreateForgotPassword(forgotPassword models.ForgotPassword) error {
	if err := fpr.db.Create(&forgotPassword).Error; err != nil {
		return err
	}
	return nil
}

// GetForgotPasswordByEmail fetches a forgot password entry by email
func (fpr *ForgotPasswordRepository) GetForgotPasswordByEmail(email string) (*models.ForgotPassword, error) {
	var forgotPassword models.ForgotPassword
	if err := fpr.db.Where("email = ?", email).First(&forgotPassword).Error; err != nil {
		return nil, err
	}
	return &forgotPassword, nil
}

// DeleteForgotPasswordByEmail deletes a forgot password entry by email
func (fpr *ForgotPasswordRepository) DeleteForgotPasswordByEmail(email string) error {
	if err := fpr.db.Where("email = ?", email).Unscoped().Delete(&models.ForgotPassword{}).Error; err != nil {
		return err
	}
	return nil
}

// **DeletionConfirmationRepository** //

type DeletionConfirmationRepository struct {
	db *gorm.DB
}

func NewDeletionConfirmationRepository() *DeletionConfirmationRepository {
	return &DeletionConfirmationRepository{database.DB}
}

type DeletionConfirmationRepositoryInterface interface {
	CreateDeletionConfirmation(deletionConfirmation models.DeletionConfirmation) error
	GetDeletionConfirmationByEmail(email string) (models.DeletionConfirmation, error)
	DeleteDeletionConfirmationByEmail(email string) error
}

// CreateDeletionConfirmation creates a new deletion confirmation entry
func (dcr *DeletionConfirmationRepository) CreateDeletionConfirmation(deletionConfirmation models.DeletionConfirmation) error {
	if err := dcr.db.Create(&deletionConfirmation).Error; err != nil {
		return err
	}
	return nil
}

// GetDeletionConfirmationByEmail fetches a deletion confirmation entry by email
func (dcr *DeletionConfirmationRepository) GetDeletionConfirmationByEmail(email string) (models.DeletionConfirmation, error) {
	var deletionConfirmation models.DeletionConfirmation
	if err := dcr.db.Where("email = ?", email).First(&deletionConfirmation).Error; err != nil {
		return deletionConfirmation, err
	}
	return deletionConfirmation, nil
}

// DeleteDeletionConfirmationByEmail deletes a deletion confirmation entry by email
func (dcr *DeletionConfirmationRepository) DeleteDeletionConfirmationByEmail(email string) error {
	if err := dcr.db.Where("email = ?", email).Unscoped().Delete(&models.DeletionConfirmation{}).Error; err != nil {
		return err
	}
	return nil
}

// **AuthProviderRepository** //

type AuthProviderRepository struct {
	db *gorm.DB
}

func NewAuthProviderRepository() *AuthProviderRepository {
	return &AuthProviderRepository{database.DB}
}

type AuthProviderRepositoryInterface interface {
	CreateAuthProvider(authProvider models.AuthProvider) error
	GetAuthProviderByID(id uint) (models.AuthProvider, error)
	GetAuthProviderByProviderKey(providerKey string) (models.AuthProvider, error)
	UpdateAuthProvider(authProvider models.AuthProvider) error
	DeleteAuthProvider(id uint) error
	DeleteAuthProviderByUserID(userID uint) error
}

// CreateAuthProvider creates a new AuthProvider record.
func (ar *AuthProviderRepository) CreateAuthProvider(authProvider models.AuthProvider) error {
	return ar.db.Create(&authProvider).Error
}

// GetAuthProviderByID retrieves an AuthProvider record by its ID.
func (ar *AuthProviderRepository) GetAuthProviderByID(id uint) (models.AuthProvider, error) {
	var authProvider models.AuthProvider
	if err := ar.db.First(&authProvider, id).Error; err != nil {
		return authProvider, err
	}
	return authProvider, nil
}

// GetAuthProviderByProviderKey retrieves an AuthProvider record by its ProviderKey.
func (ar *AuthProviderRepository) GetAuthProviderByProviderKey(providerKey string) (models.AuthProvider, error) {
	var authProvider models.AuthProvider
	if err := ar.db.Where("provider_key = ?", providerKey).First(&authProvider).Error; err != nil {
		return authProvider, err
	}
	return authProvider, nil
}

// UpdateAuthProvider updates an existing AuthProvider record.
func (ar *AuthProviderRepository) UpdateAuthProvider(authProvider models.AuthProvider) error {
	return ar.db.Save(authProvider).Error
}

// DeleteAuthProvider deletes an AuthProvider record by its ID.
func (ar *AuthProviderRepository) DeleteAuthProvider(id uint) error {
	return ar.db.Unscoped().Delete(&models.AuthProvider{}, id).Error
}

// DeleteAuthProviderByUserID deletes an AuthProvider record by its UserID.
func (ar *AuthProviderRepository) DeleteAuthProviderByUserID(userID uint) error {
	return ar.db.Where("user_id = ?", userID).Unscoped().Delete(&models.AuthProvider{}).Error
}
