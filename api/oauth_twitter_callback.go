package main

import (
	"encoding/json"
	"errors"
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

	var res twitterOAuthReponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}
	if err := res.validate(); err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	email := res.Email
	name := res.Name
	link := res.getLinkURL()
	photo := res.getImageURL()

	c, err := commenterGetByEmail("twitter", email)
	if err != nil && err != errorNoSuchCommenter {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	var commenterHex string

	if err == errorNoSuchCommenter {
		commenterHex, err = commenterNew(email, name, link, photo, "twitter", "")
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
	} else {
		if err = commenterUpdate(c.CommenterHex, email, name, link, photo, "twitter"); err != nil {
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

// response from Twitter API.
// ref: https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/user-object
type twitterOAuthReponse struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	// normal image size is 48x48.
	// ref: https://developer.twitter.com/en/docs/accounts-and-users/user-profile-images-and-banners
	ImageURL string `json:"profile_image_url_https"`
}

func (r twitterOAuthReponse) validate() error {
	if r.Email == "" {
		return errors.New("no email address returned by Twitter")
	}
	if r.Name == "" {
		return errors.New("no name returned by Twitter")
	}
	return nil
}

func (r twitterOAuthReponse) getLinkURL() string {
	if r.ScreenName == "" {
		return "undefined"
	}
	return fmt.Sprintf("https://twitter.com/%s", r.ScreenName)
}

func (r twitterOAuthReponse) getImageURL() string {
	if r.ImageURL == "" {
		return "undefined"
	}
	return r.ImageURL
}
