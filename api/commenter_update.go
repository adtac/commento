package main

import ()

func commenterUpdate(commenterHex string, email string, name string, link string, photo string, provider string) error {
	if email == "" || name == "" || link == "" || photo == "" || provider == "" {
		return errorMissingField
	}

	// See utils_sanitise.go's documentation on isHttpsUrl. This is not a URL
	// validator, just an XSS preventor.
	// TODO: reject URLs instead of malforming them.
	if link != "undefined" && !isHttpsUrl(link) {
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
