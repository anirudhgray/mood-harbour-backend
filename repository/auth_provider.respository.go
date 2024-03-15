package repository

import (
	"github.com/anirudhgray/mood-harbour-backend/infra/database"
	"github.com/anirudhgray/mood-harbour-backend/models" // Import your model package
	"gorm.io/gorm"
)

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
