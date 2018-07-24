package main

import (
	"net/http"
	"time"
)

func ownerSendResetHex(email string) error {
	if email == "" {
		return errorMissingField
	}

	if !smtpConfigured {
		return errorSmtpNotConfigured
	}

	o, err := ownerGetByEmail(email)
	if err != nil {
		if err == errorNoSuchEmail {
			// TODO: use a more random time instead.
			time.Sleep(1 * time.Second)
			return nil
		} else {
			logger.Errorf("cannot get owner by email: %v", err)
			return errorInternal
		}
	}

	resetHex, err := randomHex(32)
	if err != nil {
		return err
	}

	statement := `
		INSERT INTO
		ownerResetHexes (resetHex, ownerHex, sendDate)
		VALUES          ($1,       $2,    $3      );
	`
	_, err = db.Exec(statement, resetHex, o.OwnerHex, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert resetHex: %v", err)
		return errorInternal
	}

	err = smtpOwnerResetHex(email, o.Name, resetHex)
	if err != nil {
		return err
	}

	return nil
}

func ownerSendResetHexHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email *string `json:"email"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := ownerSendResetHex(*x.Email); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
