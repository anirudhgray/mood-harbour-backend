package models

import "gorm.io/gorm"

// MoodType represents the type of mood that a user can have.
// The mood types are: Angry, Sad, Neutral, Happy, Excited.
type MoodType int

const (
	Angry MoodType = iota + 1
	Sad
	Neutral
	Happy
	Excited
)

// Mood represents a single mood entry made by a specific user.
type Mood struct {
	gorm.Model
	UserID uint     // Foreign key to the User model
	Mood   MoodType `gorm:"not null"`
	Notes  string   `gorm:"size:255;"` // Notes are optional
}

// Attribute represents a single attribute that can be associated with a mood entry. This can be anything you might want associated with a mood entry.
type Attribute struct {
	gorm.Model
	Name string `gorm:"size:255;not null;unique"`
}

// MoodAttribute represents the attributes of a mood entry. Many to many relationship with Mood.
type MoodAttribute struct {
	gorm.Model
	MoodID      uint // Foreign key to the Mood model
	AttributeID uint // Foreign key to the Attribute model
}
