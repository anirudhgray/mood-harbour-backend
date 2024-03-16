package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"size:255;not null;unique;"`
	Name         string `gorm:"size:255;not null;"`
	ProfileImage string `gorm:"size:255;"`
	Verified     bool   `gorm:"default:false"`
	Admin        bool   `gorm:"default:false"`
	Disabled     bool   `gorm:"default:false"`
}

type DeletionConfirmation struct {
	gorm.Model
	Email     string `gorm:"unique"`
	OTP       string
	ValidTill time.Time
}

type ForgotPassword struct {
	gorm.Model
	Email     string `gorm:"unique"`
	OTP       string
	ValidTill time.Time
}

type PasswordAuth struct {
	gorm.Model
	Email    string `gorm:"size:255;not null;unique;"`
	Password string `gorm:"size:255;not null;"`
	UserID   uint   `gorm:"unique;"`
}

func (pa *PasswordAuth) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pa.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	pa.Password = string(hashedPassword)
	return nil
}

type VerificationEntry struct {
	gorm.Model
	Email string `gorm:"unique"`
	OTP   string
}

type AuthProvider struct {
	gorm.Model
	ProviderKey  string `gorm:"size:255;not null"`
	UserID       uint   // Foreign key to the User model
	ProviderName string `gorm:"size:255;not null"`
}
