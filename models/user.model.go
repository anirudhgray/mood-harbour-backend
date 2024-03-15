package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"size:255;not null;unique;"`
	Name         string `gorm:"size:255;not null;"`
	ProfileImage string `gorm:"size:255;"`
	Verified     bool   `gorm:"default:false"`
}
