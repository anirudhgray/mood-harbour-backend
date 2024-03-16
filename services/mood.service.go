package services

import (
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/repository"
)

type MoodService struct {
	moodRepo repository.MoodRepositoryInterface
}

func NewMoodService(moodRepo repository.MoodRepositoryInterface) *MoodService {
	return &MoodService{moodRepo}
}

type MoodServiceInterface interface {
	CreateMoodEntry(moodType models.MoodType, notes string, userID uint, attributes []string) (models.MoodResponse, error)
	GetUserMoodEntries(userID uint, moodType *models.MoodType, startDate, endDate *string) ([]models.MoodResponse, error)
	GetSingleUserMoodEntry(userID, moodID uint) (models.MoodResponse, error)
	DeleteMoodEntry(userID, moodID uint) error
	CreateNewAttribute(attribute string) error
	UpdateUserMoodEntry(moodID uint, moodType models.MoodType, notes string, attributes []string) error
}

// CreateMoodEntry creates a new mood entry in the database.
func (ms *MoodService) CreateMoodEntry(moodType models.MoodType, notes string, userID uint, attributes []string) (models.MoodResponse, error) {
	mood := models.Mood{
		UserID: userID,
		Mood:   moodType,
		Notes:  notes,
	}

	err := ms.moodRepo.CreateMoodEntry(&mood)
	if err != nil {
		return models.MoodResponse{}, err
	}

	for _, attr := range attributes {
		attribute := models.Attribute{Name: attr}
		err = ms.moodRepo.CreateNewAttribute(&attribute)
		if err != nil {
			return models.MoodResponse{}, err
		}

		moodAttribute := models.MoodAttribute{
			MoodID:      mood.ID,
			AttributeID: attribute.ID,
		}
		err = ms.moodRepo.CreateMoodAttributeEntry(&moodAttribute)
		if err != nil {
			return models.MoodResponse{}, err
		}
	}

	moodAttributes, err := ms.moodRepo.GetMoodAttributesByMoodID(mood.ID)
	if err != nil {
		return models.MoodResponse{}, err
	}

	// for each moodAttr get the attribute
	var attributesList []models.Attribute
	for _, moodAttr := range moodAttributes {
		attribute, err := ms.moodRepo.GetAttributeByID(moodAttr.AttributeID)
		if err != nil {
			return models.MoodResponse{}, err
		}
		attributesList = append(attributesList, attribute)
	}

	return models.MoodResponse{
		ID:         mood.ID,
		Mood:       mood,
		Attributes: attributesList,
	}, nil
}

// GetUsersMoodEntries gets all the mood entries for a specific user. Filters by mood type and date range if provided.
func (ms *MoodService) GetUserMoodEntries(userID uint, moodType *models.MoodType, startDate, endDate *string) ([]models.MoodResponse, error) {
	var moods []models.MoodResponse

	moodEntries := []models.Mood{}
	err := error(nil)
	if moodType == nil && startDate == nil && endDate == nil {
		moodEntries, err = ms.moodRepo.GetMoodsByUserID(userID)
		if err != nil {
			return nil, err
		}
	} else if moodType != nil && startDate == nil && endDate == nil {
		moodEntries, err = ms.moodRepo.GetMoodsByUserIDAndMoodType(userID, *moodType)
		if err != nil {
			return nil, err
		}
	} else if moodType == nil && startDate != nil && endDate != nil {
		moodEntries, err = ms.moodRepo.GetMoodsByUserIDAndDateRange(userID, *startDate, *endDate)
		if err != nil {
			return nil, err
		}
	} else if moodType != nil && startDate != nil && endDate != nil {
		moodEntries, err = ms.moodRepo.GetMoodsByUserIDAndMoodTypeAndDateRange(userID, *moodType, *startDate, *endDate)
		if err != nil {
			return nil, err
		}
	}

	for _, mood := range moodEntries {
		moodAttributes, err := ms.moodRepo.GetMoodAttributesByMoodID(mood.ID)
		if err != nil {
			return nil, err
		}

		// for each moodAttr get the attribute
		var attributes []models.Attribute
		for _, moodAttr := range moodAttributes {
			attribute, err := ms.moodRepo.GetAttributeByID(moodAttr.AttributeID)
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, attribute)
		}

		moods = append(moods, models.MoodResponse{
			ID:         mood.ID,
			Mood:       mood,
			Attributes: attributes,
		})
	}

	return moods, nil
}

// GetSingleUserMoodEntry gets a single mood entry for a specific user.
func (ms *MoodService) GetSingleUserMoodEntry(userID, moodID uint) (models.MoodResponse, error) {
	mood, err := ms.moodRepo.GetMoodByID(moodID)
	if err != nil {
		return models.MoodResponse{}, err
	}

	moodAttributes, err := ms.moodRepo.GetMoodAttributesByMoodID(mood.ID)
	if err != nil {
		return models.MoodResponse{}, err
	}

	// for each moodAttr get the attribute
	var attributes []models.Attribute
	for _, moodAttr := range moodAttributes {
		attribute, err := ms.moodRepo.GetAttributeByID(moodAttr.AttributeID)
		if err != nil {
			return models.MoodResponse{}, err
		}
		attributes = append(attributes, attribute)
	}

	return models.MoodResponse{
		ID:         mood.ID,
		Mood:       mood,
		Attributes: attributes,
	}, nil
}

// DeleteMoodEntry deletes a single mood entry for a specific user, as well as all associated mood attributes.
func (ms *MoodService) DeleteMoodEntry(userID, moodID uint) error {
	err := ms.moodRepo.DeleteMoodAttributeByMoodID(moodID)
	if err != nil {
		return err
	}

	err = ms.moodRepo.DeleteMood(moodID)
	if err != nil {
		return err
	}

	return nil
}

// CreateNewAttribute creates a new attribute in the database which can now be associated with mood entries.
func (ms *MoodService) CreateNewAttribute(attribute string) error {
	attr := models.Attribute{Name: attribute}
	return ms.moodRepo.CreateNewAttribute(&attr)
}

// UpdateUserMoodEntry updates a single mood entry for a specific user, as well as all associated mood attributes.
func (ms *MoodService) UpdateUserMoodEntry(moodID uint, moodType models.MoodType, notes string, attributes []string) error {
	mood, err := ms.moodRepo.GetMoodByID(moodID)
	if err != nil {
		return err
	}
	mood.Mood = moodType
	mood.Notes = notes

	err = ms.moodRepo.UpdateMoodEntry(&mood)
	if err != nil {
		return err
	}

	err = ms.moodRepo.DeleteMoodAttributeByMoodID(moodID)
	if err != nil {
		return err
	}

	for _, attr := range attributes {
		attribute := models.Attribute{Name: attr}
		err = ms.moodRepo.CreateNewAttribute(&attribute)
		if err != nil {
			return err
		}

		moodAttribute := models.MoodAttribute{
			MoodID:      mood.ID,
			AttributeID: attribute.ID,
		}
		err = ms.moodRepo.CreateMoodAttributeEntry(&moodAttribute)
		if err != nil {
			return err
		}
	}

	return nil
}
