package services

import (
	"testing"

	"github.com/anirudhgray/mood-harbour-backend/mocks"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
)

func TestMoodService_GetUserMoodEntries_ByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMoodRepo := mocks.NewMockMoodRepositoryInterface(ctrl)
	ms := NewMoodService(mockMoodRepo)
	userID := uint(1)

	expectedLength := 5
	mockMoodRepo.EXPECT().GetMoodsByUserID(userID).Return(make([]models.Mood, expectedLength), nil)
	mockMoodRepo.EXPECT().GetMoodAttributesByMoodID(gomock.Any()).Return(make([]models.MoodAttribute, 0), nil).Times(expectedLength)

	moods, err := ms.GetUserMoodEntries(userID, nil, nil, nil)
	if err != nil {
		t.Errorf("GetUserMoodEntries returned an error: %v", err)
	}
	if len(moods) != expectedLength {
		t.Errorf("Expected length: %d, got: %d", expectedLength, len(moods))
	}
}

func TestMoodService_GetUserMoodEntries_ByUserIDAndMoodType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMoodRepo := mocks.NewMockMoodRepositoryInterface(ctrl)
	ms := NewMoodService(mockMoodRepo)
	userID := uint(1)
	moodTypeHappy := models.Happy

	expectedLength := 2
	mockMoodRepo.EXPECT().GetMoodsByUserIDAndMoodType(userID, moodTypeHappy).Return(make([]models.Mood, expectedLength), nil)
	mockMoodRepo.EXPECT().GetMoodAttributesByMoodID(gomock.Any()).Return(make([]models.MoodAttribute, 0), nil).Times(expectedLength)

	moods, err := ms.GetUserMoodEntries(userID, &moodTypeHappy, nil, nil)
	if err != nil {
		t.Errorf("GetUserMoodEntries returned an error: %v", err)
	}
	if len(moods) != expectedLength {
		t.Errorf("Expected length: %d, got: %d", expectedLength, len(moods))
	}
}

func TestMoodService_GetUserMoodEntries_ByUserIDAndDateRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMoodRepo := mocks.NewMockMoodRepositoryInterface(ctrl)
	ms := NewMoodService(mockMoodRepo)
	userID := uint(1)
	startDate := "2021-01-01"
	endDate := "2021-01-31"

	expectedLength := 3
	mockMoodRepo.EXPECT().GetMoodsByUserIDAndDateRange(userID, startDate, endDate).Return(make([]models.Mood, expectedLength), nil)
	mockMoodRepo.EXPECT().GetMoodAttributesByMoodID(gomock.Any()).Return(make([]models.MoodAttribute, 0), nil).Times(expectedLength)

	moods, err := ms.GetUserMoodEntries(userID, nil, &startDate, &endDate)
	if err != nil {
		t.Errorf("GetUserMoodEntries returned an error: %v", err)
	}
	if len(moods) != expectedLength {
		t.Errorf("Expected length: %d, got: %d", expectedLength, len(moods))
	}
}

func TestMoodService_GetUserMoodEntries_ByUserIDAndMoodTypeAndDateRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMoodRepo := mocks.NewMockMoodRepositoryInterface(ctrl)
	ms := NewMoodService(mockMoodRepo)
	userID := uint(1)
	moodTypeHappy := models.Happy
	startDate := "2021-01-01"
	endDate := "2021-01-31"

	expectedLength := 1
	mockMoodRepo.EXPECT().GetMoodsByUserIDAndMoodTypeAndDateRange(userID, moodTypeHappy, startDate, endDate).Return(make([]models.Mood, expectedLength), nil)
	mockMoodRepo.EXPECT().GetMoodAttributesByMoodID(gomock.Any()).Return(make([]models.MoodAttribute, 0), nil).Times(expectedLength)

	moods, err := ms.GetUserMoodEntries(userID, &moodTypeHappy, &startDate, &endDate)
	if err != nil {
		t.Errorf("GetUserMoodEntries returned an error: %v", err)
	}
	if len(moods) != expectedLength {
		t.Errorf("Expected length: %d, got: %d", expectedLength, len(moods))
	}
}

func TestMoodService_CreateMoodEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMoodRepo := mocks.NewMockMoodRepositoryInterface(ctrl)
	ms := NewMoodService(mockMoodRepo)
	userID := uint(1)
	moodType := models.Happy
	notes := "Feeling good"
	attributes := []string{"attribute1", "attribute2"}

	// Mock the CreateMoodEntry method
	mockMoodRepo.EXPECT().CreateMoodEntry(gomock.Any()).Return(nil)

	// Mock the CreateNewAttribute and CreateMoodAttributeEntry methods for each attribute
	for _, attr := range attributes {
		attribute := &models.Attribute{Name: attr}
		mockMoodRepo.EXPECT().CreateNewAttribute(attribute).Return(nil)
		mockMoodRepo.EXPECT().CreateMoodAttributeEntry(gomock.Any()).Return(nil)
	}

	// Mock the GetMoodAttributesByMoodID and GetAttributeByID methods
	mockMoodRepo.EXPECT().GetMoodAttributesByMoodID(gomock.Any()).Return([]models.MoodAttribute{{AttributeID: 1}, {AttributeID: 2}}, nil)
	mockMoodRepo.EXPECT().GetAttributeByID(gomock.Any()).Return(models.Attribute{}, nil).Times(2)

	moodResponse, err := ms.CreateMoodEntry(moodType, notes, userID, attributes)
	if err != nil {
		t.Errorf("CreateMoodEntry returned an error: %v", err)
	}

	if moodResponse.Mood.Mood != moodType || moodResponse.Mood.Notes != notes || moodResponse.Mood.UserID != userID {
		t.Errorf("CreateMoodEntry returned incorrect mood")
	}
	if len(moodResponse.Attributes) != len(attributes) {
		t.Errorf("CreateMoodEntry returned incorrect number of attributes")
	}
}

func TestMoodService_UpdateUserMoodEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMoodRepo := mocks.NewMockMoodRepositoryInterface(ctrl)
	ms := NewMoodService(mockMoodRepo)
	moodID := uint(1)
	moodType := models.Happy
	notes := "Feeling good"
	attributes := []string{"attribute1", "attribute2"}

	// Mock the GetMoodByID method
	mockMoodRepo.EXPECT().GetMoodByID(moodID).Return(models.Mood{Model: gorm.Model{ID: moodID}, Mood: models.Angry, Notes: "Initial Notes"}, nil)

	// Mock the UpdateMoodEntry method
	mockMoodRepo.EXPECT().UpdateMoodEntry(gomock.Any()).Do(func(mood *models.Mood) {
		if mood.Mood != moodType {
			t.Errorf("Expected mood type: %v, got: %v", moodType, mood.Mood)
		}
		if mood.Notes != notes {
			t.Errorf("Expected notes: %v, got: %v", notes, mood.Notes)
		}
	}).Return(nil)

	// Mock the DeleteMoodAttributeByMoodID method
	mockMoodRepo.EXPECT().DeleteMoodAttributeByMoodID(moodID).Return(nil)

	// Mock the CreateNewAttribute and CreateMoodAttributeEntry methods for each attribute
	for _, attr := range attributes {
		attrModel := &models.Attribute{Name: attr}
		mockMoodRepo.EXPECT().CreateNewAttribute(attrModel).Return(nil)
		mockMoodRepo.EXPECT().CreateMoodAttributeEntry(gomock.Any()).Return(nil)
	}

	err := ms.UpdateUserMoodEntry(moodID, moodType, notes, attributes)

	if err != nil {
		t.Errorf("UpdateUserMoodEntry returned an error: %v", err)
	}
}
