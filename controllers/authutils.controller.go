package controllers

import (
	"net/http"
	"time"

	"github.com/GDGVIT/attendance-app-backend/infra/logger"
	"github.com/GDGVIT/attendance-app-backend/models"
	"github.com/GDGVIT/attendance-app-backend/utils/auth"
	"github.com/GDGVIT/attendance-app-backend/utils/email"
	"github.com/gin-gonic/gin"
)

// ForgotPasswordRequest handles forgot password requests by sending a mail with an OTP
func (uc *UserController) ForgotPasswordRequest(c *gin.Context) {
	useremail := c.Query("email")

	// Fetch the user by email
	user, err := uc.userRepo.GetUserByEmail(useremail)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Forgot Password mail sent."})
		return
	}

	// Check if a forgot password entry already exists for the user's email
	err = uc.forgotRepo.DeleteForgotPasswordByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Forgot Password mail sent."})
		return
	}

	// Send the forgot password email
	err = email.SendForgotPasswordMail(user.Email, user.ID, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "mail", "message": "Error in sending email."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Forgot Password mail sent."})
	logger.Infof("Forgot password request")
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

	useremail := c.Query("email")
	otp := c.Query("otp")

	// Fetch the forgot password entry by email
	forgotPasswordEntry, err := uc.forgotRepo.GetForgotPasswordByEmail(useremail)
	if err != nil {
		logger.Errorf("Error while verifying: %v", err.Error())
		c.JSON(http.StatusForbidden, gin.H{"error": "verification", "message": "Invalid verification. Please check email link again."})
		return
	}

	if forgotPasswordEntry.ValidTill.Before(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "otp-expiry", "message": "Password OTP has expired, please request forgot password again."})
		return
	}

	if forgotPasswordEntry.OTP == otp {
		// Fetch the user by email
		user, err := uc.userRepo.GetUserByEmail(useremail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user-fetch", "message": "Failed to fetch user"})
			return
		}
		pwdAuth, err := uc.passwordAuthRepo.GetPwdAuthItemByEmail(useremail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user-fetch", "message": "Failed to fetch user"})
			return
		}

		if !auth.CheckPasswordStrength(forgotPasswordInput.NewPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password-strength", "message": "Password not strong enough."})
			return
		}
		pwdAuth.Password = forgotPasswordInput.NewPassword
		pwdAuth.HashPassword()

		err = uc.userRepo.SaveUser(user)
		if err != nil {
			logger.Errorf("Save user after forgot and new: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "save-data", "message": "Failed to update password"})
			return
		}
		err = uc.passwordAuthRepo.UpdatePwdAuthItem(pwdAuth)
		if err != nil {
			logger.Errorf("Save user after forgot and new: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "save-data", "message": "Failed to update password"})
			return
		}

		email.GenericSendMail("Password Reset", "Password for your account was reset recently.", user.Email, user.Name)

		// Delete the forgot password entry
		err = uc.forgotRepo.DeleteForgotPasswordByEmail(useremail)
		if err != nil {
			logger.Errorf("Error while deleting forgot password entry: " + err.Error())
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password set successfully. Please proceed to login."})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "verification", "message": "Invalid verification. Password not updated."})
	}
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
	currentPwdAuth, err := uc.passwordAuthRepo.GetPwdAuthItemByEmail(currentUser.Email)
	if err != nil {
		logger.Errorf("Error getting password auth item: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch-data", "message": "Failed to get auth item."})
		return
	}

	if err := auth.VerifyPassword(resetPasswordInput.OldPassword, currentPwdAuth.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password", "message": "Please enter your current password correctly."})
		email.GenericSendMail("Password Reset Attempt", "Somebody attempted to change your password on Bookstore. Secure your account if this was not you.", currentUser.Email, currentUser.Name)
		return
	}

	if !auth.CheckPasswordStrength(resetPasswordInput.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password-strength", "message": "Password not strong enough."})
		return
	}
	currentPwdAuth.Password = resetPasswordInput.NewPassword
	currentPwdAuth.HashPassword()

	err = uc.userRepo.SaveUser(*currentUser)
	if err != nil {
		logger.Errorf("Update Password failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save-data", "message": "Failed to update password"})
		return
	}
	err = uc.passwordAuthRepo.UpdatePwdAuthItem(currentPwdAuth)
	if err != nil {
		logger.Errorf("Update Password failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save-data", "message": "Failed to update password"})
		return
	}

	email.GenericSendMail("Password Reset Successfully", "Your password for Bookstore was changed. Secure your account if this was not you.", currentUser.Email, currentUser.Name)
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

	// Check if a deletion confirmation record already exists for the user's email, and remove it
	err := uc.deletionRepo.DeleteDeletionConfirmationByEmail(currentUser.Email)
	if err != nil {
		logger.Errorf("Error while removing past deletion req: %v", err.Error())
		c.JSON(http.StatusForbidden, gin.H{"error": "verification", "message": "Error while removing past deletion request."})
		return
	}

	// Send deletion email
	err = email.SendDeletionMail(currentUser.Email, currentUser.ID, currentUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "mail", "message": "Error in sending email."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deletion request submitted"})
}

func (uc *UserController) DeleteAccount(c *gin.Context) {
	useremail := c.Query("email")
	otp := c.Query("otp")

	// Fetch the forgot password entry by email
	deletionEntry, err := uc.deletionRepo.GetDeletionConfirmationByEmail(useremail)
	if err != nil {
		logger.Errorf("Error while verifying deletion: %v", err.Error())
		c.JSON(http.StatusForbidden, gin.H{"error": "verification", "message": "Invalid verification. Please check email link again."})
		return
	}

	if deletionEntry.ValidTill.Before(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "otp-expiry", "message": "Password OTP has expired, please request account deletion again."})
		return
	}

	if deletionEntry.OTP == otp {
		// Fetch the user by email
		user, err := uc.userRepo.GetUserByEmail(useremail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user-fetch", "message": "Failed to fetch user"})
			return
		}

		err = uc.userRepo.DeleteUserByID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "deletion", "message": "Failed to delete user."})
			return
		}

		err = uc.passwordAuthRepo.DeletePwdAuthItemByEmail(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "deletion", "message": "Failed to delete user."})
			return
		}

		err = uc.authProviderRepo.DeleteAuthProviderByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "deletion", "message": "Failed to delete user."})
			return
		}

		email.GenericSendMail("Account Deleted", "Your account on GDSC Attendance App has been deleted.", user.Email, user.Name)

		// Delete the deletion request entry
		err = uc.deletionRepo.DeleteDeletionConfirmationByEmail(useremail)
		if err != nil {
			logger.Errorf("Error while deleting deletion entry: " + err.Error())
		}

		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully."})
		logger.Infof("Account deleted")
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "verification", "message": "Invalid verification. Account not deleted."})
	}
}
