package services

import (
	"errors"

	"github.com/GDGVIT/attendance-app-backend/infra/logger"
	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/GDGVIT/attendance-app-backend/repository"
	"github.com/GDGVIT/attendance-app-backend/utils/auth"
)

// AuthService handles the business logic for authentication
type AuthService struct {
	authProviderRepo repository.AuthProviderRepositoryInterface
	verificationRepo repository.VerificationRepositoryInterface
	forgotRepo       repository.ForgotPasswordRepositoryInterface
	deletionRepo     repository.DeletionConfirmationRepositoryInterface
	passwordAuthRepo repository.PasswordAuthRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	emailService     EmailServiceInterface
}

// NewAuthService returns a new AuthService
func NewAuthService(
	authProviderRepo repository.AuthProviderRepositoryInterface,
	verificationRepo repository.VerificationRepositoryInterface,
	forgotRepo repository.ForgotPasswordRepositoryInterface,
	deletionRepo repository.DeletionConfirmationRepositoryInterface,
	passwordAuthRepo repository.PasswordAuthRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	emailService EmailServiceInterface,
) *AuthService {
	return &AuthService{
		authProviderRepo: authProviderRepo,
		verificationRepo: verificationRepo,
		forgotRepo:       forgotRepo,
		deletionRepo:     deletionRepo,
		passwordAuthRepo: passwordAuthRepo,
		userRepo:         userRepo,
		emailService:     emailService,
	}
}

type AuthServiceInterface interface {
	RegisterUser(email, name, profileImage, password string) (models.User, error)
	LoginUser(email, password string) (string, models.User, error)
	RequestVerificationAgain(email string) error
}

func (as *AuthService) RegisterUser(email, name, profileImage, password string) (models.User, error) {
	// Check if the user already exists
	existingUser, _ := as.userRepo.GetUserByEmail(email)
	existingPwdAuth, _ := as.passwordAuthRepo.GetPwdAuthItemByEmail(email)

	var emptyPwdAuth models.PasswordAuth
	var emptyUser models.User
	if existingPwdAuth != emptyPwdAuth {
		// email.SendRegistrationMail("Account Alert", "Someone attempted to create an account using your email. If this was you, try applying for password reset in case you have lost access to your account.", existingUser.Email, existingUser.ID, existingUser.Name, false)
		return models.User{}, errors.New("User with that email address already exists!")
	}

	if !auth.CheckPasswordStrength(password) {
		return models.User{}, errors.New("Password not strong enough.")
	}

	user := models.User{Name: name, Email: email, ProfileImage: profileImage}
	pwdauth := models.PasswordAuth{Password: password, Email: email}

	// Hash the user's password
	if err := pwdauth.HashPassword(); err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		return models.User{}, err
	}

	// Create the user profile in the database IF there is no user profile yet
	if emptyUser == existingUser {
		if err := as.userRepo.CreateUser(user); err != nil {
			logger.Errorf("Failed to create user: %v", err)
			return models.User{}, err
		}
		as.emailService.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
		logger.Infof("New User Object Created.")
		u, _ := as.userRepo.GetUserByEmail(email)
		pwdauth.UserID = u.ID
	} else {
		// dead, not supporting social right now.
		pwdauth.UserID = existingUser.ID
	}
	// Create the password auth item in the database
	if err := as.passwordAuthRepo.CreatePwdAuthItem(&pwdauth); err != nil {
		logger.Errorf("Failed to create user: %v", err)
		return models.User{}, err
	}

	return user, nil

	// dead, not supporting social right now.
	// if emptyUser != existingUser {
	// 	c.JSON(http.StatusCreated, gin.H{"message": "Email-Password login method added."})
	// }
}

// LoginUser handles user login.
func (as *AuthService) LoginUser(email, password string) (string, models.User, error) {
	token, user, err := auth.LoginCheck(email, password)

	if err != nil {
		return "", models.User{}, err
	}

	if !user.Verified {
		return "", models.User{}, errors.New("Please verify your email before logging in.")
	}

	return token, user, nil
}

// RequestVerificationAgain handles resending verification email.
func (as *AuthService) RequestVerificationAgain(email string) error {
	user, err := as.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if user.Verified {
		return nil
	}

	_, err = as.verificationRepo.GetVerificationEntryByEmail(user.Email)
	if err == nil {
		err = as.verificationRepo.DeleteVerificationEntry(user.Email)
		if err != nil {
			return err
		}
	}

	// Send verification email
	err = as.emailService.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
	if err != nil {
		return err
	}

	return nil
}
