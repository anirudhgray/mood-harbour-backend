// VerificationEntryRepository.go

package repository

import (
	"github.com/GDGVIT/attendance-app-backend/infra/database"
	"github.com/GDGVIT/attendance-app-backend/infra/logger"
	"github.com/GDGVIT/attendance-app-backend/models"
	"gorm.io/gorm"
)

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
