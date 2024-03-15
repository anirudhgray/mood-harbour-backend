package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/GDGVIT/attendance-app-backend/config"
	"github.com/GDGVIT/attendance-app-backend/infra/logger"
	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/GDGVIT/attendance-app-backend/repository"
	"github.com/GDGVIT/attendance-app-backend/utils/auth"
	"github.com/GDGVIT/attendance-app-backend/utils/email"
	"github.com/GDGVIT/attendance-app-backend/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type UserController struct {
	userRepo         *repository.UserRepository
	forgotRepo       *repository.ForgotPasswordRepository
	verifRepo        *repository.VerificationEntryRepository
	deletionRepo     *repository.DeletionConfirmationRepository
	passwordAuthRepo *repository.PasswordAuthRepository
	authProviderRepo *repository.AuthProviderRepository
}

func NewUserController() *UserController {
	userRepo := repository.NewUserRepository()
	forgotRepo := repository.NewForgotPasswordRepository()
	verifRepo := repository.NewVerificationEntryRepository()
	deletionRepo := repository.NewDeletionConfirmationRepository()
	passwordAuthRepo := repository.NewPasswordAuthRepository()
	authProviderRepo := repository.NewAuthProviderRepository()
	return &UserController{userRepo, forgotRepo, verifRepo, deletionRepo, passwordAuthRepo, authProviderRepo}
}

// RegisterUser handles user registration
func (uc *UserController) RegisterUser(c *gin.Context) {
	var registerData struct {
		Email        string `json:"email"`
		Name         string `json:"name"`
		Password     string `json:"password"`
		ProfileImage string `json:"profile_image"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Improper JSON."})
		logger.Errorf("Failed to bind JSON: %v", err)
		return
	}

	// Check if the user already exists
	existingUser, _ := uc.userRepo.GetUserByEmail(registerData.Email)
	existingPwdAuth, _ := uc.passwordAuthRepo.GetPwdAuthItemByEmail(registerData.Email)

	var emptyPwdAuth models.PasswordAuth
	var emptyUser models.User
	if existingPwdAuth != emptyPwdAuth {
		email.SendRegistrationMail("Account Alert", "Someone attempted to create an account using your email. If this was you, try applying for password reset in case you have lost access to your account.", existingUser.Email, existingUser.ID, existingUser.Name, false)
		c.JSON(http.StatusBadRequest, gin.H{"message": "User with that email address already exists!", "error": "user-exists"})
		// c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		// we lie!
		return
	}

	if !auth.CheckPasswordStrength(registerData.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password-strength", "message": "Password not strong enough."})
		return
	}

	user := models.User{Name: registerData.Name, Email: registerData.Email, ProfileImage: registerData.ProfileImage}
	pwdauth := models.PasswordAuth{Password: registerData.Password, Email: registerData.Email}

	// Hash the user's password
	if err := pwdauth.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hashing", "message": "Failed to hash password"})
		logger.Errorf("Failed to hash password: %v", err)
		return
	}

	// Create the user profile in the database IF there is no user profile yet
	if emptyUser == existingUser {
		if err := uc.userRepo.CreateUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "creation-error", "message": "Failed to create user."})
			logger.Errorf("Failed to create user: %v", err)
			return
		}
		email.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
		c.JSON(http.StatusCreated, gin.H{"message": "User created. Verification email sent!"})
		logger.Infof("New User Created.")
		u, _ := uc.userRepo.GetUserByEmail(registerData.Email)
		pwdauth.UserID = u.ID
	} else {
		pwdauth.UserID = existingUser.ID
	}
	// Create the password auth item in the database
	if err := uc.passwordAuthRepo.CreatePwdAuthItem(&pwdauth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "creation-error", "message": "Failed to create user."})
		logger.Errorf("Failed to create user: %v", err)
		return
	}
	if emptyUser != existingUser {
		c.JSON(http.StatusCreated, gin.H{"message": "Email-Password login method added."})
	}
}

// Login handles user login
func (uc *UserController) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := auth.LoginCheck(loginData.Email, loginData.Password)

	if err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "credentials-error", "message": "The email or password is not correct"})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusForbidden, gin.H{"error": "unverified", "message": "Please verify your email before logging in."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (uc *UserController) RequestVerificationAgain(c *gin.Context) {
	useremail := c.Query("email")

	user, err := uc.userRepo.GetUserByEmail(useremail)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Verification email sent."})
		return
	}

	if user.Verified {
		c.JSON(http.StatusOK, gin.H{"message": "Verification email sent."})
		return
	}

	_, err = uc.verifRepo.GetVerificationEntryByEmail(user.Email)
	if err == nil {
		err = uc.verifRepo.DeleteVerificationEntry(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "deletion", "message": "Error deleting verification entry."})
			return
		}
	}

	// Send verification email
	err = email.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "mail", "message": "Error in sending email."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent to you again."})
	logger.Infof("Verification requested again")
}

// VerifyEmail takes your email and otp sent of registration to verify a user account.
func (uc *UserController) VerifyEmail(c *gin.Context) {
	email := c.Query("email")
	otp := c.Query("otp")

	// Fetch the verification entry by email
	verificationEntry, err := uc.verifRepo.GetVerificationEntryByEmail(email)
	if err != nil {
		logger.Errorf("Error while verifying: " + err.Error())
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification."})
		return
	}

	if verificationEntry.OTP == otp {
		// Verify the email by updating the user's verification status
		err = uc.userRepo.VerifyUserEmail(email)
		if err != nil {
			logger.Errorf("Error while verifying: " + err.Error())
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification."})
			return
		}

		// Delete the verification entry
		err = uc.verifRepo.DeleteVerificationEntry(email)
		if err != nil {
			logger.Errorf("Error while deleting verification entry: " + err.Error())
		}

		c.JSON(http.StatusOK, gin.H{"message": "Verified! You can now log in."})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification."})
	}
}

// GoogleLogin initiates google oauth2 flow
func (uc *UserController) GoogleLogin(c *gin.Context) {
	url := config.GoogleOAuthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

// GoogleCallback handles the callback from google oauth2
func (uc *UserController) GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	// Exchange the authorization code for an access token and ID token
	googletoken, err := config.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		// Handle the error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	// Use the 'accessToken' from the 'token' to fetch user data from the UserInfo endpoint
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(googletoken.AccessToken)
	userInfoResponse, err := http.Get(userInfoURL)
	if err != nil {
		// Handle the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info from Google"})
		return
	}
	defer userInfoResponse.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(userInfoResponse.Body).Decode(&userInfo); err != nil {
		// Handle the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info response"})
		return
	}

	// check if providerid in authprovider. If yes, get userid, generate JWT and log them in. If no, create authprovider entry, and if user entry not there then that as well.
	// get authprovider entry by providerid
	authProvider, _ := uc.authProviderRepo.GetAuthProviderByProviderKey(userInfo.ID)
	var emptyProviderEntry models.AuthProvider
	if authProvider != emptyProviderEntry { // i.e., found
		user, _ := uc.userRepo.GetUserByID(authProvider.UserID)
		jwt, err := token.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log in."})
			return
		}
		c.Redirect(http.StatusSeeOther, viper.GetString("FRONTEND_SOCIAL_REDIRECT")+"?token="+jwt)
		c.JSON(http.StatusOK, gin.H{"token": jwt, "user": user})
		return
	}

	// check is user entry exists for given social email, else create one
	user, _ := uc.userRepo.GetUserByEmail(userInfo.Email)
	var emptyUser models.User
	if user == emptyUser {
		user = models.User{Name: userInfo.Name, Email: userInfo.Email, ProfileImage: userInfo.Picture, Verified: true}
		if err := uc.userRepo.CreateUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user."})
			return
		}
	}

	user, _ = uc.userRepo.GetUserByEmail(userInfo.Email)
	// create authprovider entry
	authProvider = models.AuthProvider{ProviderName: "google", ProviderKey: userInfo.ID, UserID: user.ID}
	if err := uc.authProviderRepo.CreateAuthProvider(authProvider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create authprovider entry."})
		return
	}

	jwt, err := token.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log in."})
		return
	}

	c.Redirect(http.StatusSeeOther, viper.GetString("FRONTEND_SOCIAL_REDIRECT")+"?token="+jwt)
	c.JSON(http.StatusOK, gin.H{"token": jwt, "user": user})
}
