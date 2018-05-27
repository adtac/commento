package main

import (
	"time"
)

type commenter struct {
	CommenterHex string    `json:"commenterHex,omitempty"`
	Email        string    `json:"email,omitempty"`
	Name         string    `json:"name"`
	Link         string    `json:"link"`
	Photo        string    `json:"photo"`
	Provider     string    `json:"provider,omitempty"`
	JoinDate     time.Time `json:"joinDate,omitempty"`
}

func commenterIsProviderUser(provider string, email string) (bool, error) {
	if provider == "" || email == "" {
		return false, errorMissingField
	}

	statement := `
    SELECT EXISTS (
      SELECT 1
      FROM commenters
      WHERE email=$1 AND provider=$2
    );
  `
	row := db.QueryRow(statement, email, provider)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		logger.Errorf("error checking if provider user exists: %v", err)
		return false, errorInternal
	}

	return exists, nil
}
