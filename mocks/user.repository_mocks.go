package mocks

import (
	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/golang/mock/gomock"
)

// MockUserRepository is a mock for UserRepositoryInterface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// NewMockUserRepository creates a new mock for UserRepositoryInterface.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT methods for expected calls with return values
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks the CreateUser method.
func (m *MockUserRepository) CreateUser(user models.User) error {
	ret := m.ctrl.Call(m, "CreateUser", user)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// GetUserByEmail mocks the GetUserByEmail method.
func (m *MockUserRepository) GetUserByEmail(email string) (models.User, error) {
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	user, _ := ret[0].(models.User) // Type assertion for models.User
	err, _ := ret[1].(error)        // Type assertion for error
	return user, err
}

// GetUserByID mocks the GetUserByID method.
func (m *MockUserRepository) GetUserByID(userID uint) (models.User, error) {
	ret := m.ctrl.Call(m, "GetUserByID", userID)
	user, _ := ret[0].(models.User) // Type assertion for models.User
	err, _ := ret[1].(error)        // Type assertion for error
	return user, err
}

// VerifyUserEmail mocks the VerifyUserEmail method.
func (m *MockUserRepository) VerifyUserEmail(email string) error {
	ret := m.ctrl.Call(m, "VerifyUserEmail", email)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// SaveUser mocks the SaveUser method.
func (m *MockUserRepository) SaveUser(user models.User) error {
	ret := m.ctrl.Call(m, "SaveUser", user)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// DeleteUserByID mocks the DeleteUserByID method.
func (m *MockUserRepository) DeleteUserByID(userID uint) error {
	ret := m.ctrl.Call(m, "DeleteUserByID", userID)
	err, _ := ret[0].(error) // Type assertion for error
	return err
}

// MockUserRepositoryMockRecorder is a recorder for the MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// CreateUser mocks the CreateUser method.
func (m *MockUserRepositoryMockRecorder) CreateUser(user models.User) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "CreateUser", user)
}

// GetUserByEmail mocks the GetUserByEmail method.
func (m *MockUserRepositoryMockRecorder) GetUserByEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "GetUserByEmail", email)
}

// GetUserByID mocks the GetUserByID method.
func (m *MockUserRepositoryMockRecorder) GetUserByID(userID uint) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "GetUserByID", userID)
}

// VerifyUserEmail mocks the VerifyUserEmail method.
func (m *MockUserRepositoryMockRecorder) VerifyUserEmail(email string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "VerifyUserEmail", email)
}

// SaveUser mocks the SaveUser method.
func (m *MockUserRepositoryMockRecorder) SaveUser(user models.User) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "SaveUser", user)
}

// DeleteUserByID mocks the DeleteUserByID method.
func (m *MockUserRepositoryMockRecorder) DeleteUserByID(userID uint) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "DeleteUserByID", userID)
}
