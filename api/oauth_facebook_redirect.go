package main

import (
	"fmt"
	"net/http"
)

func facebookRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if facebookConfig == nil {
		logger.Errorf("facebook oauth access attempt without configuration")
		fmt.Fprintf(w, "error: this website has not configured facebook OAuth")
		return
	}

	commenterToken := r.FormValue("commenterToken")

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != errorNoSuchToken {
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	url := facebookConfig.AuthCodeURL(commenterToken)
	http.Redirect(w, r, url, http.StatusFound)
}
