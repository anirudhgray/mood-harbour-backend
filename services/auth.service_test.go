package services

import (
	"errors"
	"testing"

	"github.com/anirudhgray/mood-harbour-backend/mocks"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/golang/mock/gomock"
)

func TestAuthService_RequestVerificationAgain(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockVerificationRepo := mocks.NewMockVerificationEntryRepository(ctrl)
	mockEmailService := mocks.NewMockEmailService(ctrl)

	testEmail := "test@example.com"
	testName := "Test User"
	testUser := models.User{Email: testEmail, Verified: false, Name: testName, ProfileImage: "test"}

	mockUserRepo.EXPECT().GetUserByEmail(testEmail).Return(testUser, nil)
	mockVerificationRepo.EXPECT().GetVerificationEntryByEmail(testEmail).Return(nil, errors.New("not found"))
	mockEmailService.EXPECT().SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", testEmail, testName, true).Return(nil)

	as := &AuthService{
		userRepo:         mockUserRepo,
		verificationRepo: mockVerificationRepo,
		emailService:     mockEmailService,
	}

	err := as.RequestVerificationAgain(testEmail)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAuthService_VerifyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVerificationRepo := mocks.NewMockVerificationEntryRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	testEmail := "test@example.com"
	testOTP := "123456"
	testVerificationEntry := models.VerificationEntry{Email: testEmail, OTP: testOTP}

	as := &AuthService{
		verificationRepo: mockVerificationRepo,
		userRepo:         mockUserRepo,
	}

	testCases := []struct {
		name        string
		otp         string
		expectError bool
	}{
		{
			name:        "Correct OTP",
			otp:         testOTP,
			expectError: false,
		},
		{
			name:        "Incorrect OTP",
			otp:         "123458",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVerificationRepo.EXPECT().GetVerificationEntryByEmail(testEmail).Return(&testVerificationEntry, nil)
			if !tc.expectError {
				mockUserRepo.EXPECT().VerifyUserEmail(testEmail).Return(nil)
				mockVerificationRepo.EXPECT().DeleteVerificationEntry(testEmail).Return(nil)
			}

			err := as.VerifyEmail(testEmail, tc.otp)

			if tc.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
