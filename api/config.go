package main

import (
	"os"
	"strings"
	"path/filepath"
)

func parseConfig() error {
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Errorf("cannot load binary path: %v", err)
		return err
	}

	defaults := map[string]string{
		"POSTGRES": "postgres://postgres:postgres@localhost/commento?sslmode=disable",

		"PORT":   "8080",
		"ORIGIN": "",

		"CDN_PREFIX": "",

		"STATIC": binPath,

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
			logger.Errorf("missing %s environment variable", env)
			return errorMissingConfig
		}
	}

	if os.Getenv("CDN_PREFIX") == "" {
		os.Setenv("CDN_PREFIX", os.Getenv("ORIGIN"))
	}

	static := os.Getenv("STATIC")
	for strings.HasSuffix(static, "/") {
		static = static[0:len(static)-1]
	}

	file, err := os.Stat(static)
	if err != nil {
		logger.Errorf("cannot load %s: %v", static, err)
		return err
	}

	if !file.IsDir() {
		logger.Errorf("COMMENTO_STATIC=%s is not a directory", static)
		return errorNotADirectory
	}

	os.Setenv("STATIC", static)

	return nil
}
