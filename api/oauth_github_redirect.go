package main

import (
	"fmt"
	"net/http"
)

func githubRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if githubConfig == nil {
		logger.Errorf("github oauth access attempt without configuration")
		fmt.Fprintf(w, "error: this website has not configured github OAuth")
		return
	}

	commenterToken := r.FormValue("commenterToken")

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != errorNoSuchToken {
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	url := githubConfig.AuthCodeURL(commenterToken)
	http.Redirect(w, r, url, http.StatusFound)
}
