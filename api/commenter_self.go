package main

import (
	"net/http"
)

func commenterSelfHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true, "commenter": c})
}
