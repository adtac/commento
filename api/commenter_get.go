package main

import ()

func commenterGetByHex(commenterHex string) (commenter, error) {
	if commenterHex == "" {
		return commenter{}, errorMissingField
	}

	statement := `
    SELECT commenterHex, email, name, link, photo, provider, joinDate
    FROM commenters
    WHERE commenterHex = $1;
  `
	row := db.QueryRow(statement, commenterHex)

	c := commenter{}
	if err := row.Scan(&c.CommenterHex, &c.Email, &c.Name, &c.Link, &c.Photo, &c.Provider, &c.JoinDate); err != nil {
		// TODO: is this the only error?
		return commenter{}, errorNoSuchCommenter
	}

	return c, nil
}

func commenterGetByEmail(provider string, email string) (commenter, error) {
	if provider == "" || email == "" {
		return commenter{}, errorMissingField
	}

	statement := `
    SELECT commenterHex, email, name, link, photo, provider, joinDate
    FROM commenters
    WHERE email = $1 AND provider = $2;
  `
	row := db.QueryRow(statement, email, provider)

	c := commenter{}
	if err := row.Scan(&c.CommenterHex, &c.Email, &c.Name, &c.Link, &c.Photo, &c.Provider, &c.JoinDate); err != nil {
		// TODO: is this the only error?
		return commenter{}, errorNoSuchCommenter
	}

	return c, nil
}

func commenterGetByCommenterToken(commenterToken string) (commenter, error) {
	if commenterToken == "" {
		return commenter{}, errorMissingField
	}

	statement := `
    SELECT commenterHex
    FROM commenterSessions
    WHERE commenterToken = $1;
	`
	row := db.QueryRow(statement, commenterToken)

	var commenterHex string
	if err := row.Scan(&commenterHex); err != nil {
		// TODO: is the only error?
		return commenter{}, errorNoSuchToken
	}

	if commenterHex == "none" {
		return commenter{}, errorNoSuchToken
	}

	// TODO: use a join instead of two queries?
	return commenterGetByHex(commenterHex)
}
