package main

import (
	"net/http"
)

func ownerSelf(session string) (bool, owner) {
	o, err := ownerGetBySession(session)
	if err != nil {
		return false, owner{}
	}

	return true, o
}

func ownerSelfHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Session *string `json:"session"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	loggedIn, o := ownerSelf(*x.Session)

	writeBody(w, response{"success": true, "loggedIn": loggedIn, "owner": o})
}
