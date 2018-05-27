package main

import (
	"time"
)

func commenterNew(email string, name string, link string, photo string, provider string) (string, error) {
	if email == "" || name == "" || link == "" || photo == "" || provider == "" {
		return "", errorMissingField
	}

	commenterHex, err := randomHex(32)
	if err != nil {
		return "", errorInternal
	}

	statement := `
		INSERT INTO
		commenters (commenterHex, email, name, link, photo, provider, joinDate)
		VALUES     ($1,           $2,    $3,   $4,   $5,    $6,       $7      );
	`
	_, err = db.Exec(statement, commenterHex, email, name, link, photo, provider, time.Now().UTC())
	if err != nil {
		return "", errorInternal
	}

	return commenterHex, nil
}
