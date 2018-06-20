package main

import (
	"net/http"
)

func domainDelete(domain string) error {
	if domain == "" {
		return errorMissingField
	}

	statement := `
		DELETE FROM
		domains
		WHERE domain = $1;
	`
	_, err := db.Exec(statement, domain)
	if err != nil {
		return errorNoSuchDomain
	}

	statement = `
		DELETE FROM votes
		USING comments
		WHERE comments.commentHex = votes.commentHex AND comments.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf("cannot delete votes: %v", err)
		return errorInternal
	}

	statement = `
		DELETE FROM views
		WHERE views.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf("cannot delete views: %v", err)
		return errorInternal
	}

	statement = `
		DELETE FROM moderators
		WHERE moderators.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf("cannot delete domain moderators: %v", err)
		return errorInternal
	}

	statement = `
		DELETE FROM comments
		WHERE comments.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf(statement, domain)
		return errorInternal
	}

	return nil
}

func domainDeleteHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
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

	if err = domainDelete(*x.Domain); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true})
}
