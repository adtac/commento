package main

import (
	"fmt"
	"net/http"
)

func googleRedirectHandler(w http.ResponseWriter, r *http.Request) {
	session := r.FormValue("session")

	c, err := commenterGetBySession(session)
	if err != nil {
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	if c.CommenterHex != "none" {
		fmt.Fprintf(w, "error: that session is already in use\n")
		return
	}

	url := googleConfig.AuthCodeURL(session)
	http.Redirect(w, r, url, http.StatusFound)
}
