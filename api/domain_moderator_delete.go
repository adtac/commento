package main

import (
	"net/http"
)

func domainModeratorDelete(domain string, email string) error {
	if domain == "" || email == "" {
		return errorMissingConfig
	}

	statement := `
		DELETE FROM moderators
		WHERE domain=$1 AND email=$2;
	`
	_, err := db.Exec(statement, domain, email)
	if err != nil {
		logger.Errorf("cannot delete moderator: %v", err)
		return errorInternal
	}

	return nil
}

func domainModeratorDeleteHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Session *string `json:"session"`
		Domain  *string `json:"domain"`
		Email   *string `json:"email"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetBySession(*x.Session)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := stripDomain(*x.Domain)
	authorised, err := domainOwnershipVerify(domain, o.Email)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	if !authorised {
		writeBody(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err = domainModeratorDelete(domain, *x.Email); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true})
}
