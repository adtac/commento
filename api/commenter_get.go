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

func commenterGetBySession(session string) (commenter, error) {
	if session == "" {
		return commenter{}, errorMissingField
	}

	statement := `
    SELECT commenterHex
    FROM commenterSessions
    WHERE session = $1;
	`
	row := db.QueryRow(statement, session)

	var commenterHex string
	if err := row.Scan(&commenterHex); err != nil {
		// TODO: is the only error?
		return commenter{}, errorNoSuchSession
	}

	if commenterHex == "none" {
		return commenter{}, errorNoSuchSession
	}

	return commenterGetByHex(commenterHex)
}
