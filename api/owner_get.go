package main

import ()

var ownersRowColumns string = `
	owners.ownerHex,
	owners.email,
	owners.name,
	owners.confirmedEmail,
	owners.joinDate
`

func ownersRowScan(s sqlScanner, o *owner) error {
	return s.Scan(
		&o.OwnerHex,
		&o.Email,
		&o.Name,
		&o.ConfirmedEmail,
		&o.JoinDate,
	)
}

func ownerGetByEmail(email string) (owner, error) {
	if email == "" {
		return owner{}, errorMissingField
	}

	statement := `
		SELECT ` + ownersRowColumns + `
		FROM owners
		WHERE email=$1;
	`
	row := db.QueryRow(statement, email)

	var o owner
	if err := ownersRowScan(row, &o); err != nil {
		// TODO: Make sure this is actually no such email.
		return owner{}, errorNoSuchEmail
	}

	return o, nil
}

func ownerGetByOwnerToken(ownerToken string) (owner, error) {
	if ownerToken == "" {
		return owner{}, errorMissingField
	}

	statement := `
		SELECT ` + ownersRowColumns + `
		FROM owners
		WHERE owners.ownerHex IN (
			SELECT ownerSessions.ownerHex FROM ownerSessions
			WHERE ownerSessions.ownerToken = $1
		);
	`
	row := db.QueryRow(statement, ownerToken)

	var o owner
	if err := ownersRowScan(row, &o); err != nil {
		logger.Errorf("cannot scan owner: %v\n", err)
		return owner{}, errorInternal
	}

	return o, nil
}

func ownerGetByOwnerHex(ownerHex string) (owner, error) {
	if ownerHex == "" {
		return owner{}, errorMissingField
	}

	statement := `
		SELECT ` + ownersRowColumns + `
		FROM owners
		WHERE ownerHex = $1;
	`
	row := db.QueryRow(statement, ownerHex)

	var o owner
	if err := ownersRowScan(row, &o); err != nil {
		logger.Errorf("cannot scan owner: %v\n", err)
		return owner{}, errorInternal
	}

	return o, nil
}
