package main

import (
	"net/http"
)

func domainUpdate(d domain) error {
	if d.SsoProvider && d.SsoUrl == "" {
		return errorMissingField
	}

	statement := `
		UPDATE domains
		SET
			name=$2,
			state=$3,
			autoSpamFilter=$4,
			requireModeration=$5,
			requireIdentification=$6,
			moderateAllAnonymous=$7,
			emailNotificationPolicy=$8,
			commentoProvider=$9,
			googleProvider=$10,
			twitterProvider=$11,
			githubProvider=$12,
			gitlabProvider=$13,
			ssoProvider=$14,
			ssoUrl=$15,
			defaultSortPolicy=$16
		WHERE domain=$1;
	`

	_, err := db.Exec(statement,
		d.Domain,
		d.Name,
		d.State,
		d.AutoSpamFilter,
		d.RequireModeration,
		d.RequireIdentification,
		d.ModerateAllAnonymous,
		d.EmailNotificationPolicy,
		d.CommentoProvider,
		d.GoogleProvider,
		d.TwitterProvider,
		d.GithubProvider,
		d.GitlabProvider,
		d.SsoProvider,
		d.SsoUrl,
		d.DefaultSortPolicy)
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
