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
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
		Email      *string `json:"email"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	authorised, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !authorised {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err = domainModeratorDelete(domain, *x.Email); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
