package main

import ()

var googleConfigured bool
var twitterConfigured bool
var facebookConfigured bool
var githubConfigured bool
var gitlabConfigured bool

func oauthConfigure() error {
	if err := googleOauthConfigure(); err != nil {
		return err
	}

	if err := twitterOauthConfigure(); err != nil {
		return err
	}

	if err := facebookOauthConfigure(); err != nil {
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
