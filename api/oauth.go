package main

import ()

var configuredOauths []string

func oauthConfigure() error {
	configuredOauths = []string{}

	if err := googleOauthConfigure(); err != nil {
		return err
	}

	if err := twitterOauthConfigure(); err != nil {
		return err
	}

	if err := githubOauthConfigure(); err != nil {
		return err
	}

	if err := gitlabOauthConfigure(); err != nil {
		return err
	}

	return nil
}
