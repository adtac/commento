package main

import (
	"net/http"
)

func domainSsoSecretNew(domain string) (string, error) {
	if domain == "" {
		return "", errorMissingField
	}

	ssoSecret, err := randomHex(32)
	if err != nil {
		logger.Errorf("error generating SSO secret hex: %v", err)
		return "", errorInternal
	}

	statement := `
		UPDATE domains
		SET ssoSecret = $2
		WHERE domain = $1;
	`
	_, err = db.Exec(statement, domain, ssoSecret)
	if err != nil {
		logger.Errorf("cannot update ssoSecret: %v", err)
		return "", errorInternal
	}

	return ssoSecret, nil
}

func domainSsoSecretNewHandler(w http.ResponseWriter, r *http.Request) {
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

	ssoSecret, err := domainSsoSecretNew(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "ssoSecret": ssoSecret})
}
