package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     "fd",
		ClientSecret: "",
		RedirectURL:  "",                                     // Your callback route
		Scopes:       []string{"openid", "profile", "email"}, // Define the scopes you need
		Endpoint:     google.Endpoint,
	}
)

func InitialiseOAuthGoogle() {
	GoogleOAuthConfig.ClientID = viper.GetString("GOOGLE_CLIENT_ID")
	GoogleOAuthConfig.ClientSecret = viper.GetString("GOOGLE_CLIENT_SECRET")
	GoogleOAuthConfig.RedirectURL = viper.GetString("GOOGLE_REDIRECT_URI")
}
