package main

import ()

func oauthConfigure() error {
	if err := googleOauthConfigure(); err != nil {
		return err
	}

	return nil
}
