package main

import (
	"net/http"
)

func commenterSelfHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Session *string `json:"session"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetBySession(*x.Session)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true, "commenter": c})
}
