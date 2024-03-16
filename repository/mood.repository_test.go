package repository

import (
	"testing"
	"time"

	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/utils/test_utils"
	"gorm.io/gorm"
)

func TestMoodRepository_CreateMoodEntry(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.Mood{})

	// Create the MoodRepository with the test database
	mr := NewMoodRepository()
	mr.db = db

	// Create a test mood entry
	mood := models.Mood{
		UserID: 1,
		Mood:   models.Happy,
	}

	// Test CreateMoodEntry function
	if err := mr.CreateMoodEntry(&mood); err != nil {
		t.Errorf("CreateMoodEntry returned an error: %v", err)
	}

	// Test GetMoodByID function to retrieve the created mood entry
	retrievedMood, err := mr.GetMoodByID(mood.ID)
	if err != nil {
		t.Errorf("GetMoodByID returned an error: %v", err)
	}
	if retrievedMood.UserID != mood.UserID {
		t.Errorf("Expected mood UserID: %d, got: %d", mood.UserID, retrievedMood.UserID)
	}
	if retrievedMood.Mood != mood.Mood {
		t.Errorf("Expected mood Mood: %d, got: %d", mood.Mood, retrievedMood.Mood)
	}
}

func TestMoodRepository_GetMoodsByUserID(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.Mood{})

	// Create the MoodRepository with the test database
	mr := NewMoodRepository()
	mr.db = db

	// Create test mood entries
	moods := []models.Mood{
		{UserID: 1, Mood: models.Happy},
		{UserID: 1, Mood: models.Sad},
		{UserID: 2, Mood: models.Angry},
	}

	for _, mood := range moods {
		if err := mr.CreateMoodEntry(&mood); err != nil {
			t.Fatalf("Failed to create test mood entry: %v", err)
		}
	}

	testCases := []struct {
		userID         uint
		expectedLength int
	}{
		{1, 2},
		{2, 1},
		{3, 0},
	}

	for _, tc := range testCases {
		moods, err := mr.GetMoodsByUserID(tc.userID)
		if err != nil {
			t.Errorf("GetMoodsByUserID returned an error: %v", err)
		}
		if len(moods) != tc.expectedLength {
			t.Errorf("Expected length: %d, got: %d", tc.expectedLength, len(moods))
		}
	}
}

func TestMoodRepository_GetMoodsByUserIDAndDateRange(t *testing.T) {
	db, err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Migrator().DropTable(&models.Mood{})

	// Create the MoodRepository with the test database
	mr := NewMoodRepository()
	mr.db = db

	// Create test mood entries
	// Create test mood entries
	moods := []models.Mood{
		{UserID: 1, Mood: models.Happy, Model: gorm.Model{CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)}},
		{UserID: 1, Mood: models.Sad, Model: gorm.Model{CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)}},
		{UserID: 2, Mood: models.Angry, Model: gorm.Model{CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)}},
	}

	for _, mood := range moods {
		if err := mr.CreateMoodEntry(&mood); err != nil {
			t.Fatalf("Failed to create test mood entry: %v", err)
		}
	}

	testCases := []struct {
		userID         uint
		startDate      string
		endDate        string
		expectedLength int
	}{
		{1, "2022-01-01", "2022-01-03", 2},
		{2, "2022-01-01", "2022-01-04", 1},
		{1, "2022-01-04", "2022-01-05", 0},
		{1, "2022-01-01", "2022-01-02", 1},
		{1, "2022-01-03", "2022-01-03", 1},
	}

	for _, tc := range testCases {
		moods, err := mr.GetMoodsByUserIDAndDateRange(tc.userID, tc.startDate, tc.endDate)
		if err != nil {
			t.Errorf("GetMoodsByUserIDAndDateRange returned an error: %v", err)
		}
		if len(moods) != tc.expectedLength {
			t.Errorf("Expected length: %d, got: %d", tc.expectedLength, len(moods))
		}
	}
}
