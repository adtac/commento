package main

import (
	"net/http"
)

func domainDelete(domain string) error {
	if domain == "" {
		return errorMissingField
	}

	statement := `
		DELETE FROM domains
		WHERE domain = $1;
	`
	_, err := db.Exec(statement, domain)
	if err != nil {
		return errorNoSuchDomain
	}

	statement = `
		DELETE FROM views
		WHERE views.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf("cannot delete domain from views: %v", err)
		return errorInternal
	}

	statement = `
		DELETE FROM moderators
		WHERE moderators.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf("cannot delete domain from moderators: %v", err)
		return errorInternal
	}

	statement = `
		DELETE FROM ssotokens
		WHERE ssotokens.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf("cannot delete domain from ssotokens: %v", err)
		return errorInternal
	}

	// comments, votes, and pages are handled by domainClear
	if err = domainClear(domain); err != nil {
		logger.Errorf("cannot clear domain: %v", err)
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

	if err = domainDelete(*x.Domain); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
