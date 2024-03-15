package services

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"

	"github.com/anirudhgray/mood-harbour-backend/infra/logger"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/repository"
	"github.com/spf13/viper"
)

// EmailService handles sending email notifications.
type EmailService struct {
}

// NewEmailService creates a new EmailService.
func NewEmailService(userRepo repository.UserRepositoryInterface) *EmailService {
	return &EmailService{}
}

func GenerateOTP(maxDigits uint32) string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%0*d", maxDigits, bi)
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GenericEmail struct {
	Subject  string         `json:"subject"`
	From     EmailAddress   `json:"from"`
	To       []EmailAddress `json:"to"`
	Category string         `json:"category"`
	Text     string         `json:"text"`
}

type EmailServiceInterface interface {
	GenericSendMail(subject string, content string, toEmail string, userName string) error
	SendRegistrationMail(subject string, content string, toEmail string, userID uint, userName string, newUser bool) error
}

func (es *EmailService) SendRegistrationMail(subject string, content string, toEmail string, userID uint, userName string, newUser bool) error {
	otp := ""
	if newUser {
		otp = GenerateOTP(6)
		content += "http://bookstore.anrdhmshr.tech/verify?email=" + toEmail + "&otp=" + otp
	}

	err := es.GenericSendMail(subject, content, toEmail, userName)
	if err != nil {
		return err
	}

	if newUser {
		entry := models.VerificationEntry{
			Email: toEmail,
			OTP:   otp,
		}
		verifRepo := repository.NewVerificationEntryRepository()
		verifRepo.CreateVerificationEntry(entry)
	}
	return nil
}

// GenericSendMail sends a generic email.
func (es *EmailService) GenericSendMail(subject string, content string, toEmail string, userName string) error {
	url := "https://send.api.mailtrap.io/api/send"
	method := "POST"

	data := GenericEmail{
		Subject: subject,
		From: EmailAddress{
			Email: "attapp@anrdhmshr.tech",
			Name:  "Attendance App",
		},
		To: []EmailAddress{
			{
				Email: toEmail,
				Name:  userName,
			},
		},
		Category: "AttendanceApp",
		Text:     content,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Email Error: %v", err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		logger.Errorf("Email Error: %v", err)
		return err
	}

	bearer := fmt.Sprintf("Bearer %s", viper.GetString("MAILTRAP_API_TOKEN"))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Email Error: %v", err)
		return err
	}

	defer res.Body.Close()
	return nil
}
