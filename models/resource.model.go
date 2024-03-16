package models

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	CreatedBy uint   `gorm:"not null"` // Foreign key to the User model
	Title     string `gorm:"size:255;not null"`
	Content   string `gorm:"size:10000;not null"`
	URL       string `gorm:"size:255;not null"`
	External  bool   `gorm:"not null"` // True if the resource is external, false default
	AdminPost bool   `gorm:"not null"` // True if the resource is posted by an admin, false default
}

type Rating int

const (
	OneStar Rating = iota + 1
	TwoStar
	ThreeStar
	FourStar
	FiveStar
)

type Review struct {
	gorm.Model
	ResourceID uint   `gorm:"not null"` // Foreign key to the Resource model
	UserID     uint   `gorm:"not null"` // Foreign key to the User model
	Content    string `gorm:"size:10000;not null"`
	Rating     Rating `gorm:"not null"` // Rating out of 5
}

type ResourceResponse struct {
	ID       uint     `json:"id"`
	Resource Resource `json:"resource"`
	Reviews  []Review `json:"reviews"`
}
