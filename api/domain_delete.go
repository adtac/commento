package main

import ()

func domainDelete(domain string) error {
	if domain == "" {
		return errorMissingField
	}

	statement := `
		DELETE FROM
		domains
		WHERE domain = $1;
	`
	_, err := db.Exec(statement, domain)
	if err != nil {
		return errorNoSuchDomain
	}

	statement = `
		DELETE FROM votes
		USING comments
		WHERE comments.commentHex = votes.commentHex AND comments.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf("cannot delete votes: %v", err)
		return errorInternal
	}

	statement = `
		DELETE FROM views
		WHERE views.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf("cannot delete views: %v", err)
		return errorInternal
	}

	statement = `
		DELETE FROM comments
		WHERE comments.domain = $1;
	`
	_, err = db.Exec(statement, domain)
	if err != nil {
		logger.Errorf(statement, domain)
		return errorInternal
	}

	return nil
}
