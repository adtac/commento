package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"
)

func gitlabCallbackHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken := r.FormValue("state")
	code := r.FormValue("code")

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != errorNoSuchToken {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	token, err := gitlabConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	resp, err := http.Get(os.Getenv("GITLAB_URL") + "/api/v4/user?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	logger.Infof("%v", resp.StatusCode)
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
		fmt.Fprintf(w, "Error: no email address returned by Gitlab")
		return
	}

	email := user["email"].(string)

	if user["name"] == nil {
		fmt.Fprintf(w, "Error: no name returned by Gitlab")
		return
	}

	name := user["name"].(string)

	link := "undefined"
	if user["web_url"] != nil {
		link = user["web_url"].(string)
	}

	photo := "undefined"
	if user["avatar_url"] != nil {
		photo = user["avatar_url"].(string)
	}

	c, err := commenterGetByEmail("gitlab", email)
	if err != nil && err != errorNoSuchCommenter {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	var commenterHex string

	if err == errorNoSuchCommenter {
		commenterHex, err = commenterNew(email, name, link, photo, "gitlab", "")
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
	} else {
		if err = commenterUpdate(c.CommenterHex, email, name, link, photo, "gitlab"); err != nil {
			logger.Warningf("cannot update commenter: %s", err)
			// not a serious enough to exit with an error
		}

		commenterHex = c.CommenterHex
	}

	if err := commenterSessionUpdate(commenterToken, commenterHex); err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	fmt.Fprintf(w, "<html><script>window.parent.close()</script></html>")
}
