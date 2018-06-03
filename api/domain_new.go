package main

import (
	"net/http"
	"time"
)

func domainNew(ownerHex string, name string, domain string) error {
	if ownerHex == "" || name == "" || domain == "" {
		return errorMissingField
	}

	statement := `
		INSERT INTO
		domains (ownerHex, name, domain, creationDate)
		VALUES  ($1,       $2,   $3,     $4          );
	`
	_, err := db.Exec(statement, ownerHex, name, domain, time.Now().UTC())
	if err != nil {
		// TODO: Make sure this is really the error.
		return errorDomainAlreadyExists
	}

	return nil
}

func domainNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Session *string `json:"session"`
		Name    *string `json:"name"`
		Domain  *string `json:"domain"`
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

	if err = domainNew(o.OwnerHex, *x.Name, domain); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	if err = domainModeratorNew(*x.Domain, o.Email); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true, "domain": domain})
}
