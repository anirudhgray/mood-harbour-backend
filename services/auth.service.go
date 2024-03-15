package services

import (
	"errors"
	"time"

	"github.com/anirudhgray/mood-harbour-backend/infra/logger"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/repository"
	"github.com/anirudhgray/mood-harbour-backend/utils/auth"
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
	VerifyEmail(email, otp string) error
	ForgotPasswordRequest(email string) error
	SetNewPassword(email string, otp string, newPassword string) error
	ResetPassword(user models.User, oldPassword string, newPassword string) error
	RequestDeletion(user models.User) error
	DeleteAccount(email string, otp string) error
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
		as.emailService.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.Name, true)
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
	err = as.emailService.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.Name, true)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) VerifyEmail(email string, otp string) error {
	// Fetch the verification entry by email
	verificationEntry, err := as.verificationRepo.GetVerificationEntryByEmail(email)
	if err != nil {
		return err
	}

	if verificationEntry.OTP != otp {
		return errors.New("Invalid verification")
	}

	// Verify the email by updating the user's verification status
	err = as.userRepo.VerifyUserEmail(email)
	if err != nil {
		return err
	}

	// Delete the verification entry
	err = as.verificationRepo.DeleteVerificationEntry(email)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) ForgotPasswordRequest(email string) error {
	// Fetch the user by email
	user, err := as.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil
	}

	// Check if a forgot password entry already exists for the user's email
	err = as.forgotRepo.DeleteForgotPasswordByEmail(user.Email)
	if err != nil {
		return nil
	}

	// Send the forgot password email
	err = as.emailService.SendForgotPasswordMail(user.Email, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) SetNewPassword(email string, otp string, newPassword string) error {
	// Fetch the forgot password entry by email
	forgotPasswordEntry, err := as.forgotRepo.GetForgotPasswordByEmail(email)
	if err != nil {
		return err
	}

	if forgotPasswordEntry.ValidTill.Before(time.Now()) {
		return errors.New("Password OTP has expired, please request forgot password again")
	}

	if forgotPasswordEntry.OTP != otp {
		return errors.New("Invalid verification")
	}

	// Fetch the user by email
	user, err := as.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	pwdAuth, err := as.passwordAuthRepo.GetPwdAuthItemByEmail(email)
	if err != nil {
		return err
	}

	if !auth.CheckPasswordStrength(newPassword) {
		return errors.New("Password not strong enough")
	}

	pwdAuth.Password = newPassword
	pwdAuth.HashPassword()

	err = as.userRepo.SaveUser(user)
	if err != nil {
		return err
	}

	err = as.passwordAuthRepo.UpdatePwdAuthItem(pwdAuth)
	if err != nil {
		return err
	}

	as.emailService.GenericSendMail("Password Reset", "Password for your account was reset recently.", user.Email, user.Name)

	// Delete the forgot password entry
	err = as.forgotRepo.DeleteForgotPasswordByEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) ResetPassword(user models.User, oldPassword string, newPassword string) error {
	// Fetch the password auth item by email
	currentPwdAuth, err := as.passwordAuthRepo.GetPwdAuthItemByEmail(user.Email)
	if err != nil {
		return err
	}

	if err := auth.VerifyPassword(oldPassword, currentPwdAuth.Password); err != nil {
		as.emailService.GenericSendMail("Password Reset Attempt", "Somebody attempted to change your password. Secure your account if this was not you.", user.Email, user.Name)
		return errors.New("Incorrect current password")
	}

	if !auth.CheckPasswordStrength(newPassword) {
		return errors.New("Password not strong enough")
	}

	currentPwdAuth.Password = newPassword
	currentPwdAuth.HashPassword()

	err = as.userRepo.SaveUser(user)
	if err != nil {
		return err
	}

	err = as.passwordAuthRepo.UpdatePwdAuthItem(currentPwdAuth)
	if err != nil {
		return err
	}

	as.emailService.GenericSendMail("Password Reset Successfully", "Your password was changed. Secure your account if this was not you.", user.Email, user.Name)

	return nil
}

func (as *AuthService) RequestDeletion(user models.User) error {
	// Check if a deletion confirmation record already exists for the user's email, and remove it
	err := as.deletionRepo.DeleteDeletionConfirmationByEmail(user.Email)
	if err != nil {
		return err
	}

	// Send deletion email
	err = as.emailService.SendDeletionMail(user.Email, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) DeleteAccount(email string, otp string) error {
	// Fetch the deletion confirmation entry by email
	deletionEntry, err := as.deletionRepo.GetDeletionConfirmationByEmail(email)
	if err != nil {
		return err
	}

	if deletionEntry.ValidTill.Before(time.Now()) {
		return errors.New("Password OTP has expired, please request account deletion again")
	}

	if deletionEntry.OTP != otp {
		return errors.New("Invalid verification")
	}

	// Fetch the user by email
	user, err := as.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	err = as.userRepo.DeleteUserByID(user.ID)
	if err != nil {
		return err
	}

	err = as.passwordAuthRepo.DeletePwdAuthItemByEmail(user.Email)
	if err != nil {
		return err
	}

	err = as.authProviderRepo.DeleteAuthProviderByUserID(user.ID)
	if err != nil {
		return err
	}

	as.emailService.GenericSendMail("Account Deleted", "Your account on  Mood App has been deleted.", user.Email, user.Name)

	// Delete the deletion request entry
	err = as.deletionRepo.DeleteDeletionConfirmationByEmail(email)
	if err != nil {
		return err
	}

	return nil
}
