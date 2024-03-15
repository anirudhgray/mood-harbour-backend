package models

import "gorm.io/gorm"

type AuthProvider struct {
	gorm.Model
	ProviderKey  string `gorm:"size:255;not null"`
	UserID       uint   // Foreign key to the User model
	ProviderName string `gorm:"size:255;not null"`
}
