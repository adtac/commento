package main

import ()

func emailGet(em string) (email, error) {
	statement := `
		SELECT email, unsubscribeSecretHex, lastEmailNotificationDate, pendingEmails, sendReplyNotifications, sendModeratorNotifications
		FROM emails
		WHERE email = $1;
	`
	row := db.QueryRow(statement, em)

	e := email{}
	if err := row.Scan(&e.Email, &e.UnsubscribeSecretHex, &e.LastEmailNotificationDate, &e.PendingEmails, &e.SendReplyNotifications, &e.SendModeratorNotifications); err != nil {
		// TODO: is this the only error?
		return e, errorNoSuchEmail
	}

	return e, nil
}
