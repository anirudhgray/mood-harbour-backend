package mocks

import (
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/golang/mock/gomock"
)

// MockDeletionConfirmationRepository is a mock for DeletionConfirmationRepositoryInterface.
type MockDeletionConfirmationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDeletionConfirmationRepositoryMockRecorder
}

// NewMockDeletionConfirmationRepository creates a new mock for DeletionConfirmationRepositoryInterface.
func NewMockDeletionConfirmationRepository(ctrl *gomock.Controller) *MockDeletionConfirmationRepository {
	mock := &MockDeletionConfirmationRepository{ctrl: ctrl}
	mock.recorder = &MockDeletionConfirmationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT methods for expected calls with return values
func (m *MockDeletionConfirmationRepository) EXPECT() *MockDeletionConfirmationRepositoryMockRecorder {
	return m.recorder
}

// CreateDeletionConfirmation mocks the CreateDeletionConfirmation method.
func (m *MockDeletionConfirmationRepository) CreateDeletionConfirmation(deletionConfirmation models.DeletionConfirmation) error {
	ret := m.ctrl.Call(m, "CreateDeletionConfirmation", deletionConfirmation)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// GetDeletionConfirmationByEmail mocks the GetDeletionConfirmationByEmail method.
func (m *MockDeletionConfirmationRepository) GetDeletionConfirmationByEmail(email string) (models.DeletionConfirmation, error) {
	ret := m.ctrl.Call(m, "GetDeletionConfirmationByEmail", email)
	deletionConfirmation, _ := ret[0].(models.DeletionConfirmation) // Type assertion for models.DeletionConfirmation
	err, _ := ret[1].(error)                                        // Type assertion for error
	return deletionConfirmation, err
}

// DeleteDeletionConfirmationByEmail mocks the DeleteDeletionConfirmationByEmail method.
func (m *MockDeletionConfirmationRepository) DeleteDeletionConfirmationByEmail(email string) error {
	ret := m.ctrl.Call(m, "DeleteDeletionConfirmationByEmail", email)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// MockDeletionConfirmationRepositoryMockRecorder is a recorder for the MockDeletionConfirmationRepository.
type MockDeletionConfirmationRepositoryMockRecorder struct {
	mock *MockDeletionConfirmationRepository
}

// CreateDeletionConfirmation mocks the CreateDeletionConfirmation method.
func (m *MockDeletionConfirmationRepositoryMockRecorder) CreateDeletionConfirmation(deletionConfirmation models.DeletionConfirmation) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "CreateDeletionConfirmation", deletionConfirmation)
}

// GetDeletionConfirmationByEmail mocks the GetDeletionConfirmationByEmail method.
func (m *MockDeletionConfirmationRepositoryMockRecorder) GetDeletionConfirmationByEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "GetDeletionConfirmationByEmail", email)
}

// DeleteDeletionConfirmationByEmail mocks the DeleteDeletionConfirmationByEmail method.
func (m *MockDeletionConfirmationRepositoryMockRecorder) DeleteDeletionConfirmationByEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "DeleteDeletionConfirmationByEmail", email)
}
