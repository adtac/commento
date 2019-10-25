package main

import (
	"net/http"
)

func emailGet(em string) (email, error) {
	statement := `
		SELECT email, unsubscribeSecretHex, lastEmailNotificationDate, sendReplyNotifications, sendModeratorNotifications
		FROM emails
		WHERE email = $1;
	`
	row := db.QueryRow(statement, em)

	e := email{}
	if err := row.Scan(&e.Email, &e.UnsubscribeSecretHex, &e.LastEmailNotificationDate, &e.SendReplyNotifications, &e.SendModeratorNotifications); err != nil {
		// TODO: is this the only error?
		return e, errorNoSuchEmail
	}

	return e, nil
}

func emailGetByUnsubscribeSecretHex(unsubscribeSecretHex string) (email, error) {
	statement := `
		SELECT email, unsubscribeSecretHex, lastEmailNotificationDate, sendReplyNotifications, sendModeratorNotifications
		FROM emails
		WHERE unsubscribeSecretHex = $1;
	`
	row := db.QueryRow(statement, unsubscribeSecretHex)

	e := email{}
	if err := row.Scan(&e.Email, &e.UnsubscribeSecretHex, &e.LastEmailNotificationDate, &e.SendReplyNotifications, &e.SendModeratorNotifications); err != nil {
		// TODO: is this the only error?
		return e, errorNoSuchUnsubscribeSecretHex
	}

	return e, nil
}

func emailGetHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UnsubscribeSecretHex *string `json:"unsubscribeSecretHex"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	e, err := emailGetByUnsubscribeSecretHex(*x.UnsubscribeSecretHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "email": e})
}
