package mocks

import (
	"github.com/golang/mock/gomock"
)

// MockEmailService is a mock for EmailServiceInterface.
type MockEmailService struct {
	ctrl     *gomock.Controller
	recorder *MockEmailServiceMockRecorder
}

// NewMockEmailService creates a new mock for EmailServiceInterface.
func NewMockEmailService(ctrl *gomock.Controller) *MockEmailService {
	mock := &MockEmailService{ctrl: ctrl}
	mock.recorder = &MockEmailServiceMockRecorder{mock}
	return mock
}

// EXPECT methods for expected calls with return values
func (m *MockEmailService) EXPECT() *MockEmailServiceMockRecorder {
	return m.recorder
}

// GenericSendMail mocks the GenericSendMail method.
func (m *MockEmailService) GenericSendMail(subject string, content string, toEmail string, userName string) error {
	ret := m.ctrl.Call(m, "GenericSendMail", subject, content, toEmail, userName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendRegistrationMail mocks the SendRegistrationMail method.
func (m *MockEmailService) SendRegistrationMail(subject string, content string, toEmail string, userName string, newUser bool) error {
	ret := m.ctrl.Call(m, "SendRegistrationMail", subject, content, toEmail, userName, newUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendForgotPasswordMail mocks the SendForgotPasswordMail method.
func (m *MockEmailService) SendForgotPasswordMail(toEmail string, userName string) error {
	ret := m.ctrl.Call(m, "SendForgotPasswordMail", toEmail, userName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendDeletionMail mocks the SendDeletionMail method.
func (m *MockEmailService) SendDeletionMail(toEmail string, userName string) error {
	ret := m.ctrl.Call(m, "SendDeletionMail", toEmail, userName)
	ret0, _ := ret[0].(error)
	return ret0
}

// MockEmailServiceMockRecorder is a recorder for the MockEmailService.
type MockEmailServiceMockRecorder struct {
	mock *MockEmailService
}

// GenericSendMail mocks the GenericSendMail method.
func (m *MockEmailServiceMockRecorder) GenericSendMail(subject string, content string, toEmail string, userName string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "GenericSendMail", subject, content, toEmail, userName)
}

// SendRegistrationMail mocks the SendRegistrationMail method.
func (m *MockEmailServiceMockRecorder) SendRegistrationMail(subject string, content string, toEmail string, userName string, newUser bool) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "SendRegistrationMail", subject, content, toEmail, userName, newUser)
}

// SendForgotPasswordMail mocks the SendForgotPasswordMail method.
func (m *MockEmailServiceMockRecorder) SendForgotPasswordMail(toEmail string, userName string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "SendForgotPasswordMail", toEmail, userName)
}

// SendDeletionMail mocks the SendDeletionMail method.
func (m *MockEmailServiceMockRecorder) SendDeletionMail(toEmail string, userName string) *gomock.Call {
	return m.mock.ctrl.RecordCall(m.mock, "SendDeletionMail", toEmail, userName)
}
