package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

var googleConfig *oauth2.Config

func googleOauthConfigure() error {
	googleConfig = nil
	if os.Getenv("GOOGLE_KEY") == "" && os.Getenv("GOOGLE_SECRET") == "" {
		return nil
	}

	if os.Getenv("GOOGLE_KEY") == "" {
		logger.Errorf("COMMENTO_GOOGLE_KEY not configured, but COMMENTO_GOOGLE_SECRET is set")
		return errorOauthMisconfigured
	}

	if os.Getenv("GOOGLE_SECRET") == "" {
		logger.Errorf("COMMENTO_GOOGLE_SECRET not configured, but COMMENTO_GOOGLE_KEY is set")
		return errorOauthMisconfigured
	}

	logger.Infof("loading Google OAuth config")

	googleConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("ORIGIN") + "/api/oauth/google/callback",
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	googleConfigured = true

	return nil
}
