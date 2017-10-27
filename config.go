package main

import (
	"os"
	"github.com/joho/godotenv"
)

func loadConfig() error {
	// Default value for each environment variable.
	env := map[string]string{
		"COMMENTO_PORT": "8080",
	}

	// Load configuration from the environment. Final value is governed by the
	// last config file setting the variable. For example, a COMMENTO_PORT
	// value in .env.development.local will be used even if COMMENTO_PORT
	// exists in a .env.development file
	files := []string{".env.development.local", ".env.test.local", ".env.production.local", ".env.local", ".env.development", ".env.test", ".env.production", ".env"}
	for _, file := range files {
		newEnv, err := godotenv.Read(file)
		if err == nil {
			for key, value := range newEnv {
				env[key] = value
			}
		}
	}

	// TODO: Add configuration verification. This could potentially make
	// loadConfig return a non-nil error.

	for key, value := range env {
		os.Setenv(key, value)
	}

	return nil
}
