package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/anirudhgray/mood-harbour-backend/config"
	"github.com/anirudhgray/mood-harbour-backend/infra/logger"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/repository"
	"github.com/anirudhgray/mood-harbour-backend/services"
	"github.com/anirudhgray/mood-harbour-backend/utils/token"
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

	authService services.AuthServiceInterface
}

func NewUserController(authService services.AuthServiceInterface) *UserController {
	userRepo := repository.NewUserRepository()
	forgotRepo := repository.NewForgotPasswordRepository()
	verifRepo := repository.NewVerificationEntryRepository()
	deletionRepo := repository.NewDeletionConfirmationRepository()
	passwordAuthRepo := repository.NewPasswordAuthRepository()
	authProviderRepo := repository.NewAuthProviderRepository()
	return &UserController{userRepo, forgotRepo, verifRepo, deletionRepo, passwordAuthRepo, authProviderRepo, authService}
}

// RegisterUser handles user registration.
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

	createdUser, err := uc.authService.RegisterUser(registerData.Email, registerData.Name, registerData.ProfileImage, registerData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "creation-error", "message": err.Error()})
		logger.Errorf("Failed to create user: %v", err)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
	logger.Infof("New User Created.")
}

// Login handles user login.
func (uc *UserController) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := uc.authService.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credentials-error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// RequestVerificationAgain handles request to resend verification email.
func (uc *UserController) RequestVerificationAgain(c *gin.Context) {
	useremail := c.Query("email")

	err := uc.authService.RequestVerificationAgain(useremail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "deletion", "message": "Error deleting verification entry."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent to you again."})
	logger.Infof("Verification requested again")
}

func (uc *UserController) VerifyEmail(c *gin.Context) {
	email := c.Query("email")
	otp := c.Query("otp")

	err := uc.authService.VerifyEmail(email, otp)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verified! You can now log in."})
}

// not used
// GoogleLogin initiates google oauth2 flow
func (uc *UserController) GoogleLogin(c *gin.Context) {
	url := config.GoogleOAuthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

// not used
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

// ForgotPasswordRequest handles forgot password requests by sending a mail with an OTP
func (uc *UserController) ForgotPasswordRequest(c *gin.Context) {
	email := c.Query("email")

	err := uc.authService.ForgotPasswordRequest(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "mail", "message": "Error in sending email."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Forgot Password mail sent."})
}

// SetNewPassword sets a new password for the user after forgot password request
func (uc *UserController) SetNewPassword(c *gin.Context) {
	var forgotPasswordInput struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&forgotPasswordInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Improper JSON."})
		return
	}

	email := c.Query("email")
	otp := c.Query("otp")

	err := uc.authService.SetNewPassword(email, otp, forgotPasswordInput.NewPassword)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password set successfully. Please proceed to login."})
}

// ResetPasswordController handles the reset password by logged in user
func (uc *UserController) ResetPassword(c *gin.Context) {
	var resetPasswordInput struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&resetPasswordInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	err := uc.authService.ResetPassword(*currentUser, resetPasswordInput.OldPassword, resetPasswordInput.NewPassword)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// TestAuth is a test function to check if the auth middleware is working
func (uc *UserController) TestAuth(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)
	c.JSON(http.StatusOK, gin.H{"message": "Authenticated as " + currentUser.Name})
}

func (uc *UserController) RequestDeletion(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	err := uc.authService.RequestDeletion(*currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deletion request submitted"})
}

func (uc *UserController) DeleteAccount(c *gin.Context) {
	email := c.Query("email")
	otp := c.Query("otp")

	err := uc.authService.DeleteAccount(email, otp)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully."})
}
