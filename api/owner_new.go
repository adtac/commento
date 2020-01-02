package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func ownerNew(email string, name string, password string) (string, error) {
	if email == "" || name == "" || password == "" {
		return "", errorMissingField
	}

	if os.Getenv("FORBID_NEW_OWNERS") == "true" {
		return "", errorNewOwnerForbidden
	}

	if _, err := ownerGetByEmail(email); err == nil {
		return "", errorEmailAlreadyExists
	}

	if err := emailNew(email); err != nil {
		return "", errorInternal
	}

	ownerHex, err := randomHex(32)
	if err != nil {
		logger.Errorf("cannot generate ownerHex: %v", err)
		return "", errorInternal
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("cannot generate hash from password: %v\n", err)
		return "", errorInternal
	}

	statement := `
		INSERT INTO
		owners (ownerHex, email, name, passwordHash, joinDate, confirmedEmail)
		VALUES ($1,       $2,    $3,   $4,           $5,       $6            );
	`
	_, err = db.Exec(statement, ownerHex, email, name, string(passwordHash), time.Now().UTC(), !smtpConfigured)
	if err != nil {
		// TODO: Make sure `err` is actually about conflicting UNIQUE, and not some
		// other error. If it is something else, we should probably return `errorInternal`.
		return "", errorEmailAlreadyExists
	}

	if smtpConfigured {
		confirmHex, err := randomHex(32)
		if err != nil {
			logger.Errorf("cannot generate confirmHex: %v", err)
			return "", errorInternal
		}

		statement = `
			INSERT INTO
			ownerConfirmHexes (confirmHex, ownerHex, sendDate)
			VALUES            ($1,         $2,       $3      );
		`
		_, err = db.Exec(statement, confirmHex, ownerHex, time.Now().UTC())
		if err != nil {
			logger.Errorf("cannot insert confirmHex: %v\n", err)
			return "", errorInternal
		}

		if err = smtpOwnerConfirmHex(email, name, confirmHex); err != nil {
			return "", err
		}
	}

	return ownerHex, nil
}

func ownerNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Name     *string `json:"name"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if _, err := ownerNew(*x.Email, *x.Name, *x.Password); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// Errors in creating a commenter account should not hold this up.
	_, _ = commenterNew(*x.Email, *x.Name, "undefined", "undefined", "commento", *x.Password)

	bodyMarshal(w, response{"success": true, "confirmEmail": smtpConfigured})
}
