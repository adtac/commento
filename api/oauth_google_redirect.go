package main

import (
	"fmt"
	"net/http"
)

func googleRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if googleConfig == nil {
		logger.Errorf("google oauth access attempt without configuration")
		fmt.Fprintf(w, "error: this website has not configured Google OAuth")
		return
	}

	session := r.FormValue("session")

	_, err := commenterGetBySession(session)
	if err != nil && err != errorNoSuchSession {
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	url := googleConfig.AuthCodeURL(session)
	http.Redirect(w, r, url, http.StatusFound)
}
