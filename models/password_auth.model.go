package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
