package main

import (
	"fmt"
	"net/http"
)

func emailModerateHandler(w http.ResponseWriter, r *http.Request) {
	unsubscribeSecretHex := r.FormValue("unsubscribeSecretHex")
	action := r.FormValue("action")
	commentHex := r.FormValue("commentHex")
	if commentHex == "" {
		fmt.Fprintf(w, "error: invalid commentHex")
		return
	}

	statement := `
		SELECT domain, deleted
		FROM comments
		WHERE commentHex = $1;
	`
	row := db.QueryRow(statement, commentHex)

	var domain string
	var deleted bool
	if err := row.Scan(&domain, &deleted); err != nil {
		// TODO: is this the only error?
		fmt.Fprintf(w, "error: no such comment found (perhaps it has been deleted?)")
		return
	}

	if deleted {
		fmt.Fprintf(w, "error: that comment has already been deleted")
		return
	}

	e, err := emailGetByUnsubscribeSecretHex(unsubscribeSecretHex)
	if err != nil {
		fmt.Fprintf(w, "error: %v", err.Error())
		return
	}

	isModerator, err := isDomainModerator(domain, e.Email)
	if err != nil {
		logger.Errorf("error checking if %s is a moderator: %v", e.Email, err)
		fmt.Fprintf(w, "error: %v", errorInternal)
		return
	}

	if !isModerator {
		fmt.Fprintf(w, "error: you're not a moderator for that domain")
		return
	}

	// Do not use commenterGetByEmail here because we don't know which provider
	// should be used. This was poor design on multiple fronts on my part, but
	// let's deal with that later. For now, it suffices to match the
	// deleter/approver with any account owned by the same email.
	statement = `
		SELECT commenterHex
		FROM commenters
		WHERE email = $1;
	`
	row = db.QueryRow(statement, e.Email)

	var commenterHex string
	if err = row.Scan(&commenterHex); err != nil {
		logger.Errorf("cannot retrieve commenterHex by email %q: %v", e.Email, err)
		fmt.Fprintf(w, "error: %v", errorInternal)
		return
	}

	switch action {
	case "approve":
		err = commentApprove(commentHex)
	case "delete":
		err = commentDelete(commentHex, commenterHex)
	default:
		err = errorInvalidAction
	}

	if err != nil {
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	fmt.Fprintf(w, "comment successfully %sd", action)
}
