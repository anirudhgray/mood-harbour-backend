// ForgotPasswordRepository.go

package repository

import (
	"github.com/GDGVIT/attendance-app-backend/infra/database"
	"github.com/GDGVIT/attendance-app-backend/models"
	"gorm.io/gorm"
)

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
