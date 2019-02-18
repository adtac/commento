package main

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func commenterNew(email string, name string, link string, photo string, provider string, password string) (string, error) {
	if email == "" || name == "" || link == "" || photo == "" || provider == "" {
		return "", errorMissingField
	}

	if provider == "commento" && password == "" {
		return "", errorMissingField
	}

	// See utils_sanitise.go's documentation on isHttpsUrl. This is not a URL
	// validator, just an XSS preventor.
	// TODO: reject URLs instead of malforming them.
	if link != "undefined" && !isHttpsUrl(link) {
		link = "https://" + link
	}

	if _, err := commenterGetByEmail(provider, email); err == nil {
		return "", errorEmailAlreadyExists
	}

	if err := emailNew(email); err != nil {
		return "", errorInternal
	}

	commenterHex, err := randomHex(32)
	if err != nil {
		return "", errorInternal
	}

	passwordHash := []byte{}
	if password != "" {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			logger.Errorf("cannot generate hash from password: %v\n", err)
			return "", errorInternal
		}
	}

	statement := `
		INSERT INTO
		commenters (commenterHex, email, name, link, photo, provider, passwordHash, joinDate)
		VALUES     ($1,           $2,    $3,   $4,   $5,    $6,       $7,           $8      );
	`
	_, err = db.Exec(statement, commenterHex, email, name, link, photo, provider, string(passwordHash), time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert commenter: %v", err)
		return "", errorInternal
	}

	return commenterHex, nil
}

func commenterNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Name     *string `json:"name"`
		Website  *string `json:"website"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// TODO: add gravatar?
	// TODO: email confirmation if provider = commento?
	// TODO: email confirmation if provider = commento?
	if *x.Website == "" {
		*x.Website = "undefined"
	}

	if _, err := commenterNew(*x.Email, *x.Name, *x.Website, "undefined", "commento", *x.Password); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "confirmEmail": smtpConfigured})
}
