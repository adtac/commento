package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func commenterLogin(email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", errorMissingField
	}

	statement := `
		SELECT commenterHex, passwordHash
		FROM commenters
		WHERE email = $1 AND provider = 'commento';
	`
	row := db.QueryRow(statement, email)

	var commenterHex string
	var passwordHash string
	if err := row.Scan(&commenterHex, &passwordHash); err != nil {
		return "", errorInvalidEmailPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		// TODO: is this the only possible error?
		return "", errorInvalidEmailPassword
	}

	session, err := randomHex(32)
	if err != nil {
		logger.Errorf("cannot create session hex: %v", err)
		return "", errorInternal
	}

	statement = `
		INSERT INTO
		commenterSessions (session, commenterHex, creationDate)
		VALUES            ($1,      $2,           $3          );
	`
	_, err = db.Exec(statement, session, commenterHex, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert session token: %v\n", err)
		return "", errorInternal
	}

	return session, nil
}

func commenterLoginHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	session, err := commenterLogin(*x.Email, *x.Password)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true, "session": session})
}
