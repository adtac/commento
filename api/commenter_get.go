package main

import ()

var commentersRowColumns string = `
	commenters.commenterHex,
	commenters.email,
	commenters.name,
	commenters.link,
	commenters.photo,
	commenters.provider,
	commenters.joinDate
`

func commentersRowScan(s sqlScanner, c *commenter) error {
	return s.Scan(
		&c.CommenterHex,
		&c.Email,
		&c.Name,
		&c.Link,
		&c.Photo,
		&c.Provider,
		&c.JoinDate,
	)
}

func commenterGetByHex(commenterHex string) (commenter, error) {
	if commenterHex == "" {
		return commenter{}, errorMissingField
	}

	statement := `
		SELECT ` + commentersRowColumns + `
		FROM commenters
		WHERE commenterHex = $1;
	`
	row := db.QueryRow(statement, commenterHex)

	var c commenter
	if err := commentersRowScan(row, &c); err != nil {
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
		SELECT ` + commentersRowColumns + `
		FROM commenters
		WHERE email = $1 AND provider = $2;
	`
	row := db.QueryRow(statement, email, provider)

	var c commenter
	if err := commentersRowScan(row, &c); err != nil {
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
		SELECT ` + commentersRowColumns + `
		FROM commenterSessions
		JOIN commenters ON commenterSessions.commenterHex = commenters.commenterHex
		WHERE commenterToken = $1;
	`
	row := db.QueryRow(statement, commenterToken)

	var c commenter
	if err := commentersRowScan(row, &c); err != nil {
		// TODO: is this the only error?
		return commenter{}, errorNoSuchToken
	}

	if c.CommenterHex == "none" {
		return commenter{}, errorNoSuchToken
	}

	return c, nil
}
