package main

import (
	"github.com/gomodule/oauth1/oauth"
	"os"
	"sync"
)

type twitterOauthState struct {
	CommenterToken string
	Cred           *oauth.Credentials
}

var twitterClient *oauth.Client
var twitterCredMapLock sync.RWMutex
var twitterCredMap map[string]twitterOauthState

func twitterOauthConfigure() error {
	twitterClient = nil
	if os.Getenv("TWITTER_KEY") == "" && os.Getenv("TWITTER_SECRET") == "" {
		return nil
	}

	if os.Getenv("TWITTER_KEY") == "" {
		logger.Errorf("COMMENTO_TWITTER_KEY not configured, but COMMENTO_TWITTER_SECRET is set")
		return errorOauthMisconfigured
	}

	if os.Getenv("TWITTER_SECRET") == "" {
		logger.Errorf("COMMENTO_TWITTER_SECRET not configured, but COMMENTO_TWITTER_KEY is set")
		return errorOauthMisconfigured
	}

	logger.Infof("loading twitter OAuth config")

	twitterClient = &oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
		Credentials: oauth.Credentials{
			Token:  os.Getenv("TWITTER_KEY"),
			Secret: os.Getenv("TWITTER_SECRET"),
		},
	}

	twitterCredMap = make(map[string]twitterOauthState, 1e3)

	twitterConfigured = true

	return nil
}
