package repository

import (
	"github.com/anirudhgray/mood-harbour-backend/infra/database"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"gorm.io/gorm"
)

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
