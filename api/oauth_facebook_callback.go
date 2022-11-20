package main

import (
	"fmt"
	fb "github.com/huandu/facebook/v2"
	"github.com/koron/go-dproxy"
	"golang.org/x/oauth2"
	"net/http"
)

const fbVersion = "v15.0"

func facebookCallbackHandler(w http.ResponseWriter, r *http.Request) {

	commenterToken := r.FormValue("state")
	code := r.FormValue("code")

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != errorNoSuchToken {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	token, err := facebookConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	client := facebookConfig.Client(oauth2.NoContext, token)
	session := &fb.Session{
		Version:    fbVersion,
		HttpClient: client,
	}

	resp, err := session.Get("/me?fields=id,name,email,picture,link", nil)

	if resp["id"] == nil {
		fmt.Fprintf(w, "Error: no user id returned by Facebook")
		return
	}

	if resp["email"] == nil {
		fmt.Fprintf(w, "Error: no email returned by Facebook")
		return
	}
	email := resp["email"].(string)

	if resp["name"] == nil {
		fmt.Fprintf(w, "Error: no name returned by Facebook")
		return
	}
	name := resp["name"].(string)

	link := "undefined"
	if resp["link"] != nil {
		link = resp["link"].(string)
	}

	photo := "undefined"
	pic, err := dproxy.New(resp).M("picture").M("data").M("url").String()
	if err == nil {
		photo = pic
	}

	c, err := commenterGetByEmail("facebook", email)
	if err != nil && err != errorNoSuchCommenter {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	var commenterHex string

	if err == errorNoSuchCommenter {
		commenterHex, err = commenterNew(email, name, link, photo, "facebook", "")
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
	} else {
		if err = commenterUpdate(c.CommenterHex, email, name, link, photo, "facebook"); err != nil {
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
