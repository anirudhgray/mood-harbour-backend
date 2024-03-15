package mocks

import (
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/golang/mock/gomock"
)

// MockForgotPasswordRepository is a mock for ForgotPasswordRepositoryInterface.
type MockForgotPasswordRepository struct {
	ctrl     *gomock.Controller
	recorder *MockForgotPasswordRepositoryMockRecorder
}

// NewMockForgotPasswordRepository creates a new mock for ForgotPasswordRepositoryInterface.
func NewMockForgotPasswordRepository(ctrl *gomock.Controller) *MockForgotPasswordRepository {
	mock := &MockForgotPasswordRepository{ctrl: ctrl}
	mock.recorder = &MockForgotPasswordRepositoryMockRecorder{mock}
	return mock
}

// EXPECT methods for expected calls with return values
func (m *MockForgotPasswordRepository) EXPECT() *MockForgotPasswordRepositoryMockRecorder {
	return m.recorder
}

// CreateForgotPassword mocks the CreateForgotPassword method.
func (m *MockForgotPasswordRepository) CreateForgotPassword(forgotPassword models.ForgotPassword) error {
	ret := m.ctrl.Call(m, "CreateForgotPassword", forgotPassword)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// GetForgotPasswordByEmail mocks the GetForgotPasswordByEmail method.
func (m *MockForgotPasswordRepository) GetForgotPasswordByEmail(email string) (*models.ForgotPassword, error) {
	ret := m.ctrl.Call(m, "GetForgotPasswordByEmail", email)
	forgotPassword, _ := ret[0].(*models.ForgotPassword) // Type assertion for *models.ForgotPassword
	err, _ := ret[1].(error)                             // Type assertion for error
	return forgotPassword, err
}

// DeleteForgotPasswordByEmail mocks the DeleteForgotPasswordByEmail method.
func (m *MockForgotPasswordRepository) DeleteForgotPasswordByEmail(email string) error {
	ret := m.ctrl.Call(m, "DeleteForgotPasswordByEmail", email)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// MockForgotPasswordRepositoryMockRecorder is a recorder for the MockForgotPasswordRepository.
type MockForgotPasswordRepositoryMockRecorder struct {
	mock *MockForgotPasswordRepository
}

// CreateForgotPassword mocks the CreateForgotPassword method.
func (m *MockForgotPasswordRepositoryMockRecorder) CreateForgotPassword(forgotPassword models.ForgotPassword) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "CreateForgotPassword", forgotPassword)
}

// GetForgotPasswordByEmail mocks the GetForgotPasswordByEmail method.
func (m *MockForgotPasswordRepositoryMockRecorder) GetForgotPasswordByEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "GetForgotPasswordByEmail", email)
}

// DeleteForgotPasswordByEmail mocks the DeleteForgotPasswordByEmail method.
func (m *MockForgotPasswordRepositoryMockRecorder) DeleteForgotPasswordByEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "DeleteForgotPasswordByEmail", email)
}
