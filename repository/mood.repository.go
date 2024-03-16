package repository

import (
	"github.com/anirudhgray/mood-harbour-backend/models"
	"gorm.io/gorm"
)

type MoodRepository struct {
	db *gorm.DB
}

func NewMoodRepository(db *gorm.DB) *MoodRepository {
	return &MoodRepository{db}
}

// MoodRepositoryInterface is the interface for the MoodRepository.
type MoodRepositoryInterface interface {
	CreateMoodEntry(mood *models.Mood) error
	GetMoodsByUserID(userID uint) ([]models.Mood, error)
	GetMoodsByUserIDAndMoodType(userID uint, moodType models.MoodType) ([]models.Mood, error)
	GetMoodsByUserIDAndDateRange(userID uint, startDate, endDate string) ([]models.Mood, error)
	GetMoodByID(moodID uint) (models.Mood, error)
	GetMoodAttributesByMoodID(moodID uint) ([]models.MoodAttribute, error)
	DeleteMood(moodID uint) error
	CreateNewAttribute(attribute *models.Attribute) error
	CreateMoodAttributeEntry(moodAttribute *models.MoodAttribute) error
}

// CreateMoodEntry creates a new mood entry in the database.
func (mr *MoodRepository) CreateMoodEntry(mood *models.Mood) error {
	return mr.db.Create(mood).Error
}

// GetMoodsByUserID gets all the mood entries for a specific user.
func (mr *MoodRepository) GetMoodsByUserID(userID uint) ([]models.Mood, error) {
	var moods []models.Mood
	err := mr.db.Where("user_id = ?", userID).Find(&moods).Error
	return moods, err
}

// GetMoodsByUserIDAndMoodType gets all the mood entries for a specific user and mood type.
func (mr *MoodRepository) GetMoodsByUserIDAndMoodType(userID uint, moodType models.MoodType) ([]models.Mood, error) {
	var moods []models.Mood
	err := mr.db.Where("user_id = ? AND mood = ?", userID, moodType).Find(&moods).Error
	return moods, err
}

// GetMoodsByUserIDAndDateRange gets all the mood entries for a specific user within a date range.
func (mr *MoodRepository) GetMoodsByUserIDAndDateRange(userID uint, startDate, endDate string) ([]models.Mood, error) {
	var moods []models.Mood
	err := mr.db.Where("user_id = ? AND DATE(created_at) BETWEEN ? AND ?", userID, startDate, endDate).Find(&moods).Error
	return moods, err
}

// GetMoodByID gets a specific mood entry by its ID.
func (mr *MoodRepository) GetMoodByID(moodID uint) (models.Mood, error) {
	var mood models.Mood
	err := mr.db.First(&mood, moodID).Error
	return mood, err
}

// GetMoodAttributesByMoodID gets all the attributes associated with a specific mood entry.
func (mr *MoodRepository) GetMoodAttributesByMoodID(moodID uint) ([]models.MoodAttribute, error) {
	var moodAttributes []models.MoodAttribute
	err := mr.db.Where("mood_id = ?", moodID).Find(&moodAttributes).Error
	return moodAttributes, err
}

// DeleteMood deletes a specific mood entry by its ID.
func (mr *MoodRepository) DeleteMood(moodID uint) error {
	return mr.db.Delete(&models.Mood{}, moodID).Error
}

// CreateNewAttribute creates a new attribute in the database which can now be associated with mood entries.
func (mr *MoodRepository) CreateNewAttribute(attribute *models.Attribute) error {
	return mr.db.Create(attribute).Error
}

// CreateMoodAttributeEntry creates a new entry in the MoodAttribute table which associates an attribute with a mood entry.
func (mr *MoodRepository) CreateMoodAttributeEntry(moodAttribute *models.MoodAttribute) error {
	return mr.db.Create(moodAttribute).Error
}
