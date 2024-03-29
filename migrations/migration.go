package migrations

import (
	"github.com/anirudhgray/mood-harbour-backend/infra/database"
	"github.com/anirudhgray/mood-harbour-backend/models"
)

// Migrate Add list of model add for migrations
func Migrate() {
	var migrationModels = []interface{}{
		&models.User{},
		&models.VerificationEntry{},
		&models.ForgotPassword{},
		&models.DeletionConfirmation{},
		&models.PasswordAuth{},
		&models.AuthProvider{},
		&models.Mood{},
		&models.MoodAttribute{},
		&models.Attribute{},
		&models.Resource{},
		&models.Review{},
	}
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}

	// Remove the 'Password' field from the 'users' table
	// database.DB.Migrator().DropColumn(&models.User{}, "password")
}
