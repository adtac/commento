package main

import (
	"fmt"
	"net/http"
	"os"
)

func ownerConfirmHex(confirmHex string) error {
	if confirmHex == "" {
		return errorMissingField
	}

	statement := `
		UPDATE owners
		SET confirmedEmail=true
		WHERE ownerHex IN (
			SELECT ownerHex FROM ownerConfirmHexes
			WHERE confirmHex=$1
		);
	`
	res, err := db.Exec(statement, confirmHex)
	if err != nil {
		logger.Errorf("cannot mark user's confirmedEmail as true: %v\n", err)
		return errorInternal
	}

	count, err := res.RowsAffected()
	if err != nil {
		logger.Errorf("cannot count rows affected: %v\n", err)
		return errorInternal
	}

	if count == 0 {
		return errorNoSuchConfirmationToken
	}

	statement = `
		DELETE FROM ownerConfirmHexes
		WHERE confirmHex=$1;
	`
	_, err = db.Exec(statement, confirmHex)
	if err != nil {
		logger.Warningf("cannot remove confirmation token: %v\n", err)
		// Don't return an error because this is not critical.
	}

	return nil
}

func ownerConfirmHexHandler(w http.ResponseWriter, r *http.Request) {
	if confirmHex := r.FormValue("confirmHex"); confirmHex != "" {
		if err := ownerConfirmHex(confirmHex); err == nil {
			http.Redirect(w, r, fmt.Sprintf("%s/login?confirmed=true", os.Getenv("ORIGIN")), http.StatusTemporaryRedirect)
			return
		}
	}

	// TODO: include error message in the URL
	http.Redirect(w, r, fmt.Sprintf("%s/login?confirmed=false", os.Getenv("ORIGIN")), http.StatusTemporaryRedirect)
}
