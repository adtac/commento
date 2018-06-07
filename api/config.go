package main

import (
	"os"
)

func parseConfig() error {
	defaults := map[string]string{
		"POSTGRES": "postgres://postgres:postgres@localhost/commento?sslmode=disable",

		"PORT":   "8080",
		"ORIGIN": "",

		"CDN_PREFIX": "",

		"SMTP_USERNAME":     "",
		"SMTP_PASSWORD":     "",
		"SMTP_HOST":         "",
		"SMTP_FROM_ADDRESS": "",

		"OAUTH_GOOGLE_KEY":    "",
		"OAUTH_GOOGLE_SECRET": "",
	}

	for key, value := range defaults {
		if os.Getenv("COMMENTO_" + key) == "" {
			os.Setenv(key, value)
		} else {
			os.Setenv(key, os.Getenv("COMMENTO_" + key))
		}
	}

	// Mandatory config parameters
	for _, env := range []string{"POSTGRES", "PORT", "ORIGIN"} {
		if os.Getenv(env) == "" {
			logger.Fatalf("missing %s environment variable", env)
			return errorMissingConfig
		}
	}

	if os.Getenv("CDN_PREFIX") == "" {
		os.Setenv("CDN_PREFIX", os.Getenv("ORIGIN"))
	}

	return nil
}
