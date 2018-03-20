package main

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func loadConfig() error {
	env := map[string]string{
		"COMMENTO_DATABASE_FILE": "commento.sqlite3",
		"COMMENTO_PORT":          "8080",
	}

	// Configuration precedence (highest to lowest):
	//   Environment variables (overwrites everything below)
	//   .env
	//   .env.production
	//   .env.test
	//   .env.development
	//   .env.local
	//   .env.production.local
	//   .env.test.local
	//   .env.development.local
	files := []string{".env", ".env.production", ".env.test", ".env.development", ".env.local", ".env.production.local", ".env.test.local", ".env.development.local"}
	for _, file := range files {
		newEnv, err := godotenv.Read(file)
		if err == nil {
			for key, value := range newEnv {
				key = strings.TrimSpace(key)
				value = strings.TrimSpace(value)
				env[key] = value
			}
		}
	}

	// TODO: Add configuration verification. This could potentially make
	// loadConfig return a non-nil error.

	for key, value := range env {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return nil
}
