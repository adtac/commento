package main

import ()

func ownerGetByEmail(email string) (owner, error) {
	if email == "" {
		return owner{}, errorMissingField
	}

	statement := `
    SELECT ownerHex, email, name, confirmedEmail, joinDate
    FROM owners
    WHERE email=$1;
  `
	row := db.QueryRow(statement, email)

	var o owner
	if err := row.Scan(&o.OwnerHex, &o.Email, &o.Name, &o.ConfirmedEmail, &o.JoinDate); err != nil {
		// TODO: Make sure this is actually no such email.
		return owner{}, errorNoSuchEmail
	}

	return o, nil
}

func ownerGetBySession(session string) (owner, error) {
	if session == "" {
		return owner{}, errorMissingField
	}

	statement := `
    SELECT ownerHex, email, name, confirmedEmail, joinDate
		FROM owners
		WHERE email IN (
			SELECT email FROM ownerSessions
			WHERE session=$1
		);
	`
	row := db.QueryRow(statement, session)

	var o owner
	if err := row.Scan(&o.OwnerHex, &o.Email, &o.Name, &o.ConfirmedEmail, &o.JoinDate); err != nil {
		logger.Errorf("cannot scan owner: %v\n", err)
		return owner{}, errorInternal
	}

	return o, nil
}
