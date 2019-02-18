package main

import (
	"fmt"
	"net/http"
)

func emailModerateHandler(w http.ResponseWriter, r *http.Request) {
	unsubscribeSecretHex := r.FormValue("unsubscribeSecretHex")
	e, err := emailGetByUnsubscribeSecretHex(unsubscribeSecretHex)
	if err != nil {
		fmt.Fprintf(w, "error: %v", err.Error())
		return
	}

	action := r.FormValue("action")
	if action != "delete" && action != "approve" {
		fmt.Fprintf(w, "error: invalid action")
		return
	}

	commentHex := r.FormValue("commentHex")
	if commentHex == "" {
		fmt.Fprintf(w, "error: invalid commentHex")
		return
	}

	statement := `
		SELECT domain
		FROM comments
		WHERE commentHex = $1;
	`
	row := db.QueryRow(statement, commentHex)

	var domain string
	if err = row.Scan(&domain); err != nil {
		// TODO: is this the only error?
		fmt.Fprintf(w, "error: no such comment found (perhaps it has been deleted?)")
		return
	}

	isModerator, err := isDomainModerator(domain, e.Email)
	if err != nil {
		logger.Errorf("error checking if %s is a moderator: %v", e.Email, err)
		fmt.Fprintf(w, "error checking if %s is a moderator: %v", e.Email, err)
		return
	}

	if !isModerator {
		fmt.Fprintf(w, "error: you're not a moderator for that domain")
		return
	}

	if action == "approve" {
		err = commentApprove(commentHex)
	} else {
		err = commentDelete(commentHex)
	}

	if err != nil {
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	fmt.Fprintf(w, "comment successfully %sd", action)
}
