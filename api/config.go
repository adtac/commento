package main

import (
	"os"
	"path/filepath"
	"strings"
)

func configParse() error {
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Errorf("cannot load binary path: %v", err)
		return err
	}

	defaults := map[string]string{
		"CONFIG_FILE": "",

		"POSTGRES": "postgres://postgres:postgres@localhost/commento?sslmode=disable",

		"BIND_ADDRESS": "127.0.0.1",
		"PORT":         "8080",
		"ORIGIN":       "",

		"CDN_PREFIX": "",

		"FORBID_NEW_OWNERS": "false",

		"STATIC": binPath,

		"GZIP_STATIC": "false",

		"SMTP_USERNAME":     "",
		"SMTP_PASSWORD":     "",
		"SMTP_HOST":         "",
		"SMTP_PORT":         "",
		"SMTP_FROM_ADDRESS": "",

		"GOOGLE_KEY":    "",
		"GOOGLE_SECRET": "",
	}

	for key, value := range defaults {
		if os.Getenv("COMMENTO_"+key) == "" {
			os.Setenv(key, value)
		} else {
			os.Setenv(key, os.Getenv("COMMENTO_"+key))
		}
	}

	if os.Getenv("CONFIG_FILE") != "" {
		if err := configFileLoad(os.Getenv("CONFIG_FILE")); err != nil {
			return err
		}
	}

	// Mandatory config parameters
	for _, env := range []string{"POSTGRES", "PORT", "ORIGIN", "FORBID_NEW_OWNERS"} {
		if os.Getenv(env) == "" {
			logger.Errorf("missing COMMENTO_%s environment variable", env)
			return errorMissingConfig
		}
	}

	if os.Getenv("CDN_PREFIX") == "" {
		os.Setenv("CDN_PREFIX", os.Getenv("ORIGIN"))
	}

	if os.Getenv("FORBID_NEW_OWNERS") != "true" && os.Getenv("FORBID_NEW_OWNERS") != "false" {
		logger.Errorf("COMMENTO_FORBID_NEW_OWNERS neither 'true' nor 'false'")
		return errorInvalidConfigValue
	}

	static := os.Getenv("STATIC")
	for strings.HasSuffix(static, "/") {
		static = static[0 : len(static)-1]
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
