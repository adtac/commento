package main

import ()

var configuredOauths []string

func oauthConfigure() error {
	configuredOauths = []string{}

	if err := googleOauthConfigure(); err != nil {
		return err
	}

	return nil
}
