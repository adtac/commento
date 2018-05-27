package main

import (
	"net/http"
	"time"
)

func domainModeratorNew(domain string, email string) error {
	if domain == "" || email == "" {
		return errorMissingField
	}

	statement := `
		INSERT INTO
		moderators (domain, email, addDate)
		VALUES     ($1,     $2,    $3     );
	`
	_, err := db.Exec(statement, domain, email, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert new moderator: %v", err)
		return errorInternal
	}

	return nil
}

func domainModeratorNewHandler(w http.ResponseWriter, r *http.Request) {
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
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		writeBody(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err = domainModeratorNew(domain, *x.Email); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true})
}
