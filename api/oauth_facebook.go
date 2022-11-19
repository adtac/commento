package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"os"
)

var facebookConfig *oauth2.Config

func facebookOauthConfigure() error {
	facebookConfig = nil
	if os.Getenv("FACEBOOK_KEY") == "" && os.Getenv("FACEBOOK_SECRET") == "" {
		return nil
	}

	if os.Getenv("FACEBOOK_KEY") == "" {
		logger.Errorf("COMMENTO_FACEBOOK_KEY not configured, but COMMENTO_FACEBOOK_SECRET is set")
		return errorOauthMisconfigured
	}

	if os.Getenv("FACEBOOK_SECRET") == "" {
		logger.Errorf("COMMENTO_FACEBOOK_SECRET not configured, but COMMENTO_FACEBOOK_KEY is set")
		return errorOauthMisconfigured
	}

	logger.Infof("loading facebook OAuth config")

	facebookConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("ORIGIN") + "/api/oauth/facebook/callback",
		ClientID:     os.Getenv("FACEBOOK_KEY"),
		ClientSecret: os.Getenv("FACEBOOK_SECRET"),
		Scopes: []string{
			"public_profile",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
	facebookConfig.Endpoint.AuthURL  = "https://www.facebook.com/dialog/oauth"
        facebookConfig.Endpoint.TokenURL = "https://graph.facebook.com/oauth/access_token"

        facebookConfigured = true

	return nil
}
