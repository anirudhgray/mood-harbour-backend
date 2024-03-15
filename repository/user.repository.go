package repository

import (
	"github.com/GDGVIT/attendance-app-backend/infra/database" // Import your custom database package
	"github.com/GDGVIT/attendance-app-backend/infra/logger"
	"github.com/GDGVIT/attendance-app-backend/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{database.DB}
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
