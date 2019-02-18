package main

import (
	"time"
)

func emailNew(email string) error {
	unsubscribeSecretHex, err := randomHex(32)
	if err != nil {
		return errorInternal
	}

	statement := `
		INSERT INTO
		emails (email, unsubscribeSecretHex, lastEmailNotificationDate)
		VALUES ($1,    $2,                   $3                       )
		ON CONFLICT DO NOTHING;
	`
	_, err = db.Exec(statement, email, unsubscribeSecretHex, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert email into emails: %v", err)
		return errorInternal
	}

	return nil
}
