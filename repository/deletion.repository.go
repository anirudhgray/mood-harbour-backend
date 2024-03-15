package repository

import (
	"github.com/GDGVIT/attendance-app-backend/infra/database"
	"github.com/GDGVIT/attendance-app-backend/models"
	"gorm.io/gorm"
)

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
