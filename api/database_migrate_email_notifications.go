package main

import ()

func migrateEmails() error {
	statement := `
		SELECT commenters.email
		FROM commenters
		UNION
		SELECT owners.email
		FROM owners
		UNION
		SELECT moderators.email
		FROM moderators;
	`
	rows, err := db.Query(statement)
	if err != nil {
		logger.Errorf("cannot get comments: %v", err)
		return errorDatabaseMigration
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err = rows.Scan(&email); err != nil {
			logger.Errorf("cannot get email from tables during migration: %v", err)
			return errorDatabaseMigration
		}

		if err = emailNew(email); err != nil {
			logger.Errorf("cannot insert email during migration: %v", err)
			return errorDatabaseMigration
		}
	}

	return nil
}
