package main

import (
	"net/http"
	"time"
)

func forgot(email string, entity string) error {
	if email == "" {
		return errorMissingField
	}

	if entity != "owner" && entity != "commenter" {
		return errorInvalidEntity
	}

	if !smtpConfigured {
		return errorSmtpNotConfigured
	}

	var hex string
	var name string
	if entity == "owner" {
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
		hex = o.OwnerHex
		name = o.Name
	} else {
		c, err := commenterGetByEmail("commento", email)
		if err != nil {
			if err == errorNoSuchEmail {
				// TODO: use a more random time instead.
				time.Sleep(1 * time.Second)
				return nil
			} else {
				logger.Errorf("cannot get commenter by email: %v", err)
				return errorInternal
			}
		}
		hex = c.CommenterHex
		name = c.Name
	}

	resetHex, err := randomHex(32)
	if err != nil {
		return err
	}

	var statement string

	statement = `
		INSERT INTO
		resetHexes (resetHex, hex, entity, sendDate)
		VALUES     ($1,       $2,  $3,     $4      );
	`
	_, err = db.Exec(statement, resetHex, hex, entity, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert resetHex: %v", err)
		return errorInternal
	}

	err = smtpResetHex(email, name, resetHex)
	if err != nil {
		return err
	}

	return nil
}

func forgotHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email  *string `json:"email"`
		Entity *string `json:"entity"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := forgot(*x.Email, *x.Entity); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
