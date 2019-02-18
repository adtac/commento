package main

import (
	"net/http"
	"time"
)

func domainModeratorNew(domain string, email string) error {
	if domain == "" || email == "" {
		return errorMissingField
	}

	if err := emailNew(email); err != nil {
		logger.Errorf("cannot create email when creating moderator: %v", err)
		return errorInternal
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
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err = domainModeratorNew(domain, *x.Email); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
