package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func twitterCallbackHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("oauth_token")
	verifier := r.FormValue("oauth_verifier")

	twitterCredMapLock.RLock()
	s, ok := twitterCredMap[token]
	twitterCredMapLock.RUnlock()

	commenterToken := s.CommenterToken

	if !ok {
		fmt.Fprintf(w, "no such token/verifier combination found")
		return
	}

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != errorNoSuchToken {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	x, _, err := twitterClient.RequestToken(nil, s.Cred, verifier)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	twitterCredMapLock.Lock()
	delete(twitterCredMap, token)
	twitterCredMapLock.Unlock()

	resp, err := twitterClient.Get(nil, x, "https://api.twitter.com/1.1/account/verify_credentials.json", url.Values{"include_email": {"true"}})
	if err != nil {
		fmt.Fprintf(w, "Error getting email: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "Error: status %d: %s\n", resp.StatusCode, msg)
		return
	}

	var res map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	if res["email"] == nil {
		fmt.Fprintf(w, "Error: no email address returned by Twitter")
		return
	}

	email := res["email"].(string)

	if res["name"] == nil {
		fmt.Fprintf(w, "Error: no name returned by Twitter")
		return
	}

	name := res["name"].(string)

	link := "undefined"
	photo := "undefined"
	if res["handle"] != nil {
		handle := res["screen_name"].(string)
		link = "https://twitter.com/" + handle
		photo = "https://twitter.com/" + handle + "/profile_image"
	}

	c, err := commenterGetByEmail("twitter", email)
	if err != nil && err != errorNoSuchCommenter {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	var commenterHex string

	// TODO: in case of returning users, update the information we have on record?
	if err == errorNoSuchCommenter {
		commenterHex, err = commenterNew(email, name, link, photo, "twitter", "")
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
