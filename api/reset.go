package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func reset(resetHex string, password string) (string, error) {
	if resetHex == "" || password == "" {
		return "", errorMissingField
	}

	statement := `
		SELECT hex, entity
		FROM resetHexes
		WHERE resetHex = $1;
	`
	row := db.QueryRow(statement, resetHex)

	var hex string
	var entity string
	if err := row.Scan(&hex, &entity); err != nil {
		// TODO: is this the only error?
		return "", errorNoSuchResetToken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("cannot generate hash from password: %v\n", err)
		return "", errorInternal
	}

	if entity == "owner" {
		statement = `
			UPDATE owners SET passwordHash = $1
			WHERE ownerHex = $2;
		`
	} else {
		statement = `
			UPDATE commenters SET passwordHash = $1
			WHERE commenterHex = $2;
		`
	}

	_, err = db.Exec(statement, string(passwordHash), hex)
	if err != nil {
		logger.Errorf("cannot change %s's password: %v\n", entity, err)
		return "", errorInternal
	}

	statement = `
		DELETE FROM resetHexes
		WHERE resetHex = $1;
	`
	_, err = db.Exec(statement, resetHex)
	if err != nil {
		logger.Warningf("cannot remove resetHex: %v\n", err)
	}

	return entity, nil
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ResetHex *string `json:"resetHex"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	entity, err := reset(*x.ResetHex, *x.Password)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "entity": entity})
}
