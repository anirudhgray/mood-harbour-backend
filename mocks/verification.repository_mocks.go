package mocks

import (
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/golang/mock/gomock"
)

// MockVerificationEntryRepository is a mock for VerificationRepositoryInterface.
type MockVerificationEntryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockVerificationEntryRepositoryMockRecorder
}

// NewMockVerificationEntryRepository creates a new mock for VerificationRepositoryInterface.
func NewMockVerificationEntryRepository(ctrl *gomock.Controller) *MockVerificationEntryRepository {
	mock := &MockVerificationEntryRepository{ctrl: ctrl}
	mock.recorder = &MockVerificationEntryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT methods for expected calls with return values
func (m *MockVerificationEntryRepository) EXPECT() *MockVerificationEntryRepositoryMockRecorder {
	return m.recorder
}

// CreateVerificationEntry mocks the CreateVerificationEntry method.
func (m *MockVerificationEntryRepository) CreateVerificationEntry(verificationEntry models.VerificationEntry) error {
	ret := m.ctrl.Call(m, "CreateVerificationEntry", verificationEntry)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// GetVerificationEntryByEmail mocks the GetVerificationEntryByEmail method.
func (m *MockVerificationEntryRepository) GetVerificationEntryByEmail(email string) (*models.VerificationEntry, error) {
	ret := m.ctrl.Call(m, "GetVerificationEntryByEmail", email)
	verificationEntry, _ := ret[0].(*models.VerificationEntry) // Type assertion for *models.VerificationEntry
	err, _ := ret[1].(error)                                   // Type assertion for error
	return verificationEntry, err
}

// DeleteVerificationEntry mocks the DeleteVerificationEntry method.
func (m *MockVerificationEntryRepository) DeleteVerificationEntry(email string) error {
	ret := m.ctrl.Call(m, "DeleteVerificationEntry", email)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// MockVerificationEntryRepositoryMockRecorder is a recorder for the MockVerificationEntryRepository.
type MockVerificationEntryRepositoryMockRecorder struct {
	mock *MockVerificationEntryRepository
}

// CreateVerificationEntry mocks the CreateVerificationEntry method.
func (m *MockVerificationEntryRepositoryMockRecorder) CreateVerificationEntry(verificationEntry models.VerificationEntry) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "CreateVerificationEntry", verificationEntry)
}

// GetVerificationEntryByEmail mocks the GetVerificationEntryByEmail method.
func (m *MockVerificationEntryRepositoryMockRecorder) GetVerificationEntryByEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "GetVerificationEntryByEmail", email)
}

// DeleteVerificationEntry mocks the DeleteVerificationEntry method.
func (m *MockVerificationEntryRepositoryMockRecorder) DeleteVerificationEntry(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "DeleteVerificationEntry", email)
}
