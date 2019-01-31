package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken := r.FormValue("state")
	code := r.FormValue("code")

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != errorNoSuchToken {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	token, err := googleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", errorCannotReadResponse.Error())
		return
	}

	user := make(map[string]interface{})
	if err := json.Unmarshal(contents, &user); err != nil {
		fmt.Fprintf(w, "Error: %s", errorInternal.Error())
		return
	}

	if user["email"] == nil {
		fmt.Fprintf(w, "Error: no email address returned by Github")
		return
	}

	email := user["email"].(string)

	c, err := commenterGetByEmail("google", email)
	if err != nil && err != errorNoSuchCommenter {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	var commenterHex string

	// TODO: in case of returning users, update the information we have on record?
	if err == errorNoSuchCommenter {
		var link string
		if val, ok := user["link"]; ok {
			link = val.(string)
		} else {
			link = "undefined"
		}

		commenterHex, err = commenterNew(email, user["name"].(string), link, user["picture"].(string), "google", "")
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
	} else {
		commenterHex = c.CommenterHex
	}

	if err := commenterSessionUpdate(commenterToken, commenterHex); err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	fmt.Fprintf(w, "<html><script>window.parent.close()</script></html>")
}
