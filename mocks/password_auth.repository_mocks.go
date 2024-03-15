package mocks

import (
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/golang/mock/gomock"
)

// MockPasswordAuthRepository is a mock for PasswordAuthRepositoryInterface.
type MockPasswordAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordAuthRepositoryMockRecorder
}

// NewMockPasswordAuthRepository creates a new mock for PasswordAuthRepositoryInterface.
func NewMockPasswordAuthRepository(ctrl *gomock.Controller) *MockPasswordAuthRepository {
	mock := &MockPasswordAuthRepository{ctrl: ctrl}
	mock.recorder = &MockPasswordAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT methods for expected calls with return values
func (m *MockPasswordAuthRepository) EXPECT() *MockPasswordAuthRepositoryMockRecorder {
	return m.recorder
}

// CreatePwdAuthItem mocks the CreatePwdAuthItem method.
func (m *MockPasswordAuthRepository) CreatePwdAuthItem(passwordAuth *models.PasswordAuth) error {
	ret := m.ctrl.Call(m, "CreatePwdAuthItem", passwordAuth)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// GetPwdAuthItemByEmail mocks the GetPwdAuthItemByEmail method.
func (m *MockPasswordAuthRepository) GetPwdAuthItemByEmail(email string) (models.PasswordAuth, error) {
	ret := m.ctrl.Call(m, "GetPwdAuthItemByEmail", email)
	passwordAuth, _ := ret[0].(models.PasswordAuth) // Type assertion for models.PasswordAuth
	err, _ := ret[1].(error)                        // Type assertion for error
	return passwordAuth, err
}

// UpdatePwdAuthItem mocks the UpdatePwdAuthItem method.
func (m *MockPasswordAuthRepository) UpdatePwdAuthItem(passwordAuth models.PasswordAuth) error {
	ret := m.ctrl.Call(m, "UpdatePwdAuthItem", passwordAuth)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// DeletePwdAuthItem mocks the DeletePwdAuthItem method.
func (m *MockPasswordAuthRepository) DeletePwdAuthItem(id uint) error {
	ret := m.ctrl.Call(m, "DeletePwdAuthItem", id)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// DeletePwdAuthItemByEmail mocks the DeletePwdAuthItemByEmail method.
func (m *MockPasswordAuthRepository) DeletePwdAuthItemByEmail(email string) error {
	ret := m.ctrl.Call(m, "DeletePwdAuthItemByEmail", email)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// MockPasswordAuthRepositoryMockRecorder is a recorder for the MockPasswordAuthRepository.
type MockPasswordAuthRepositoryMockRecorder struct {
	mock *MockPasswordAuthRepository
}

// CreatePwdAuthItem mocks the CreatePwdAuthItem method.
func (m *MockPasswordAuthRepositoryMockRecorder) CreatePwdAuthItem(passwordAuth *models.PasswordAuth) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "CreatePwdAuthItem", passwordAuth)
}

// GetPwdAuthItemByEmail mocks the GetPwdAuthItemByEmail method.
func (m *MockPasswordAuthRepositoryMockRecorder) GetPwdAuthItemByEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "GetPwdAuthItemByEmail", email)
}

// UpdatePwdAuthItem mocks the UpdatePwdAuthItem method.
func (m *MockPasswordAuthRepositoryMockRecorder) UpdatePwdAuthItem(passwordAuth models.PasswordAuth) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "UpdatePwdAuthItem", passwordAuth)
}

// DeletePwdAuthItem mocks the DeletePwdAuthItem method.
func (m *MockPasswordAuthRepositoryMockRecorder) DeletePwdAuthItem(id uint) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "DeletePwdAuthItem", id)
}

// DeletePwdAuthItemByEmail mocks the DeletePwdAuthItemByEmail method.
func (m *MockPasswordAuthRepositoryMockRecorder) DeletePwdAuthItemByEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "DeletePwdAuthItemByEmail", email)
}
