package main

import ()

func commentDomainGet(commentHex string) (string, error) {
	if commentHex == "" {
		return "", errorMissingField
	}

	statement := `
    SELECT domain
		FROM comments
		WHERE commentHex = $1;
	`
	row := db.QueryRow(statement, commentHex)

	var domain string
	var err error
	if err = row.Scan(&domain); err != nil {
		return "", errorNoSuchDomain
	}

	return domain, nil
}
