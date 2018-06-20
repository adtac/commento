package main

import (
	"net/http"
)

func ownerSelfHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err == errorNoSuchToken {
		writeBody(w, response{"success": true, "loggedIn": false})
		return
	}

	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true, "loggedIn": true, "owner": o})
}
