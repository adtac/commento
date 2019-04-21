package main

import (
	"time"
)

type ssoPayload struct {
	Domain string `json:"domain"`
	Token  string `json:"token"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Link   string `json:"link"`
	Photo  string `json:"photo"`
}

func ssoTokenNew(domain string, commenterToken string) (string, error) {
	token, err := randomHex(32)
	if err != nil {
		logger.Errorf("error generating SSO token hex: %v", err)
		return "", errorInternal
	}

	statement := `
		INSERT INTO
		ssoTokens (token, domain, commenterToken, creationDate)
		VALUES    ($1,    $2,     $3,             $4          );
	`
	_, err = db.Exec(statement, token, domain, commenterToken, time.Now().UTC())
	if err != nil {
		logger.Errorf("error inserting SSO token: %v", err)
		return "", errorInternal
	}

	return token, nil
}

func ssoTokenExtract(token string) (string, string, error) {
	statement := `
		SELECT domain, commenterToken
		FROM ssoTokens
		WHERE token = $1;
	`
	row := db.QueryRow(statement, token)

	var domain string
	var commenterToken string
	if err := row.Scan(&domain, &commenterToken); err != nil {
		return "", "", errorNoSuchToken
	}

	statement = `
		DELETE FROM ssoTokens
		WHERE token = $1;
	`
	if _, err := db.Exec(statement, token); err != nil {
		logger.Errorf("cannot delete SSO token after usage: %v", err)
		return "", "", errorInternal
	}

	return domain, commenterToken, nil
}
