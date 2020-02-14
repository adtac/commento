package main

import (
	"net/http"
)

func commenterUpdate(commenterHex string, email string, name string, link string, photo string, provider string) error {
	if email == "" || name == "" || photo == "" || provider == "" {
		return errorMissingField
	}

	// See utils_sanitise.go's documentation on isHttpsUrl. This is not a URL
	// validator, just an XSS preventor.
	// TODO: reject URLs instead of malforming them.
	if link == "" {
		link = "undefined"
	} else if link != "undefined" && !isHttpsUrl(link) {
		link = "https://" + link
	}

	statement := `
		UPDATE commenters
		SET email = $3, name = $4, link = $5, photo = $6
		WHERE commenterHex = $1 and provider = $2;
	`
	_, err := db.Exec(statement, commenterHex, provider, email, name, link, photo)
	if err != nil {
		logger.Errorf("cannot update commenter: %v", err)
		return errorInternal
	}

	return nil
}

func commenterUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		Name           *string `json:"name"`
		Email          *string `json:"email"`
		Link           *string `json:"link"`
		Photo          *string `json:"photo"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if c.Provider != "commento" {
		bodyMarshal(w, response{"success": false, "message": errorCannotUpdateOauthProfile.Error()})
		return
	}

	*x.Email = c.Email

	if err = commenterUpdate(c.CommenterHex, *x.Email, *x.Name, *x.Link, *x.Photo, c.Provider); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
