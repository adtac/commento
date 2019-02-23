package main

import (
	"fmt"
	"net/http"
	"os"
)

func twitterRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if twitterClient == nil {
		logger.Errorf("twitter oauth access attempt without configuration")
		fmt.Fprintf(w, "error: this website has not configured twitter OAuth")
		return
	}

	commenterToken := r.FormValue("commenterToken")

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != errorNoSuchToken {
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	cred, err := twitterClient.RequestTemporaryCredentials(nil, os.Getenv("ORIGIN")+"/api/oauth/twitter/callback", nil)
	if err != nil {
		logger.Errorf("cannot get temporary twitter credentials: %v", err)
		fmt.Fprintf(w, "error: %v", errorInternal.Error())
		return
	}

	twitterCredMapLock.Lock()
	twitterCredMap[cred.Token] = twitterOauthState{
		CommenterToken: commenterToken,
		Cred:           cred,
	}
	twitterCredMapLock.Unlock()

	http.Redirect(w, r, twitterClient.AuthorizationURL(cred, nil), http.StatusFound)
}
