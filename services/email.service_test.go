package services

import (
	"testing"

	"github.com/anirudhgray/mood-harbour-backend/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailService_GenericSendMail(t *testing.T) {
	// Create a new instance of the gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	emailService := NewEmailService(mockUserRepo)

	// Define your test data
	subject := "Test Subject"
	content := "Test Content"
	toEmail := "anirudh04mishra@gmail.com"
	userName := "Test User"

	// Call the GenericSendMail method
	err := emailService.GenericSendMail(subject, content, toEmail, userName)

	// Assert that no error occurred during the execution
	assert.NoError(t, err)
}

// TODO need to add tests
