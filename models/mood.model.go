package models

import "gorm.io/gorm"

// MoodType represents the type of mood that a user can have.
// The mood types are: Angry, Sad, Neutral, Happy, Excited.
type MoodType int
type AttributeQuantity int

const (
	Angry MoodType = iota + 1
	Sad
	Neutral
	Happy
	Excited
)

const (
	Low AttributeQuantity = iota + 1
	Medium
	High
)

// Mood represents a single mood entry made by a specific user.
type Mood struct {
	gorm.Model
	UserID uint     `gorm:"not null"` // Foreign key to the User model
	Mood   MoodType `gorm:"not null"`
	Notes  string   `gorm:"size:255;"` // Notes are optional
}

// Attribute represents a single attribute that can be associated with a mood entry. This can be anything you might want associated with a mood entry.
type Attribute struct {
	gorm.Model
	Name      string            `gorm:"size:255;not null"`
	Quantity  AttributeQuantity `gorm:"not null"`
	CreatedBy uint              `gorm:"not null"` // Foreign key to the User model
}

// MoodAttribute represents the attributes of a mood entry. Many to many relationship with Mood.
type MoodAttribute struct {
	gorm.Model
	MoodID      uint `gorm:"primaryKey;not null"` // Foreign key to the Mood model
	AttributeID uint `gorm:"primaryKey;not null"` // Foreign key to the Attribute model
}

// MoodResponse represents the response body for a single mood entry.
type MoodResponse struct {
	ID         uint        `json:"id"`
	Mood       Mood        `json:"mood"`
	Attributes []Attribute `json:"attributes"`
}
