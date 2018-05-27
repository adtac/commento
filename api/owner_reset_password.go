package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ownerResetPassword(resetHex string, password string) error {
	if resetHex == "" || password == "" {
		return errorMissingField
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("cannot generate hash from password: %v\n", err)
		return errorInternal
	}

	statement := `
		UPDATE owners SET passwordHash=$1
		WHERE email IN (
			SELECT email FROM ownerResetHexes
			WHERE resetHex=$2
		);
	`
	res, err := db.Exec(statement, string(passwordHash), resetHex)
	if err != nil {
		logger.Errorf("cannot change user's password: %v\n", err)
		return errorInternal
	}

	count, err := res.RowsAffected()
	if err != nil {
		logger.Errorf("cannot count rows affected: %v\n", err)
		return errorInternal
	}

	if count == 0 {
		return errorNoSuchResetToken
	}

	statement = `
		DELETE FROM ownerResetHexes
    WHERE resetHex=$1;
	`
	_, err = db.Exec(statement, resetHex)
	if err != nil {
		logger.Warningf("cannot remove reset token: %v\n", err)
	}

	return nil
}

func ownerResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ResetHex *string `json:"resetHex"`
		Password *string `json:"password"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := ownerResetPassword(*x.ResetHex, *x.Password); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true})
}
