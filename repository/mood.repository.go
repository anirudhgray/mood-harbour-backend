package repository

import (
	"github.com/anirudhgray/mood-harbour-backend/infra/database"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"gorm.io/gorm"
)

type MoodRepository struct {
	db *gorm.DB
}

func NewMoodRepository() *MoodRepository {
	return &MoodRepository{database.DB}
}

// MoodRepositoryInterface is the interface for the MoodRepository.
type MoodRepositoryInterface interface {
	CreateMoodEntry(mood *models.Mood) error
	GetMoodsByUserID(userID uint) ([]models.Mood, error)
	GetMoodsByUserIDAndMoodType(userID uint, moodType models.MoodType) ([]models.Mood, error)
	GetMoodsByUserIDAndDateRange(userID uint, startDate, endDate string) ([]models.Mood, error)
	GetMoodsByUserIDAndMoodTypeAndDateRange(userID uint, moodType models.MoodType, startDate, endDate string) ([]models.Mood, error)
	GetMoodByID(moodID uint) (models.Mood, error)
	GetMoodAttributesByMoodID(moodID uint) ([]models.MoodAttribute, error)
	DeleteMood(moodID uint) error
	CreateNewAttribute(attribute *models.Attribute) error
	CreateMoodAttributeEntry(moodAttribute *models.MoodAttribute) error
	GetAttributeByID(attributeID uint) (models.Attribute, error)
	DeleteMoodAttributeByMoodID(moodID uint) error
	UpdateMoodEntry(mood *models.Mood) error
	GetAttributes(userID uint) ([]models.Attribute, error)
	GetMoodsByUserIDAndOrderedByCreatedAt(userID uint) ([]models.Mood, error)
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

// GetMoodsByUserIDAndMoodTypeAndDateRange gets all the mood entries for a specific user and mood type within a date range.
func (mr *MoodRepository) GetMoodsByUserIDAndMoodTypeAndDateRange(userID uint, moodType models.MoodType, startDate, endDate string) ([]models.Mood, error) {
	var moods []models.Mood
	err := mr.db.Where("user_id = ? AND mood = ? AND DATE(created_at) BETWEEN ? AND ?", userID, moodType, startDate, endDate).Find(&moods).Error
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

// GetAttributeByID gets a specific attribute by its ID.
func (mr *MoodRepository) GetAttributeByID(attributeID uint) (models.Attribute, error) {
	var attribute models.Attribute
	err := mr.db.First(&attribute, attributeID).Error
	return attribute, err
}

// DeleteMoodAttributeByMoodID deletes all the attributes associated with a specific mood entry
func (mr *MoodRepository) DeleteMoodAttributeByMoodID(moodID uint) error {
	return mr.db.Where("mood_id = ?", moodID).Delete(&models.MoodAttribute{}).Error
}

// UpdateMoodEntry updates a specific mood entry in the database.
func (mr *MoodRepository) UpdateMoodEntry(mood *models.Mood) error {
	return mr.db.Save(mood).Error
}

// GetAttributes gets all the attributes in the database.
func (mr *MoodRepository) GetAttributes(userID uint) ([]models.Attribute, error) {
	var attributes []models.Attribute
	err := mr.db.Where("created_by = ?", userID).Find(&attributes).Error
	return attributes, err
}

// GetMoodsByUserIDAndOrderedByCreatedAt gets all the mood entries for a specific user ordered by created_at.
func (mr *MoodRepository) GetMoodsByUserIDAndOrderedByCreatedAt(userID uint) ([]models.Mood, error) {
	var moods []models.Mood
	err := mr.db.Where("user_id = ?", userID).Order("created_at desc").Find(&moods).Error
	return moods, err
}
