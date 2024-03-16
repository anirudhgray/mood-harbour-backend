package migrations

import (
	"github.com/anirudhgray/mood-harbour-backend/infra/database"
	"github.com/anirudhgray/mood-harbour-backend/models"
)

// Migrate Add list of model add for migrations
// TODO later separate migration each models
func Migrate() {
	var migrationModels = []interface{}{
		&models.User{},
		&models.VerificationEntry{},
		&models.ForgotPassword{},
		&models.DeletionConfirmation{},
		&models.PasswordAuth{},
		&models.AuthProvider{},
	}
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}

	// Remove the 'Password' field from the 'users' table
	// database.DB.Migrator().DropColumn(&models.User{}, "password")
}
