package main

import (
	"net/http"
)

func domainUpdate(d domain) error {
	statement := `
		UPDATE domains
    SET name=$2, state=$3, autoSpamFilter=$4, requireModeration=$5, requireIdentification=$6, moderateAllAnonymous=$7
		WHERE domain=$1;
	`

	_, err := db.Exec(statement, d.Domain, d.Name, d.State, d.AutoSpamFilter, d.RequireModeration, d.RequireIdentification, d.ModerateAllAnonymous)
	if err != nil {
		logger.Errorf("cannot update non-moderators: %v", err)
		return errorInternal
	}

	return nil
}

func domainUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		D          *domain `json:"domain"`
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

	domain := domainStrip((*x.D).Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err = domainUpdate(*x.D); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
