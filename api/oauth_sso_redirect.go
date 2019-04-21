package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
)

func ssoRedirectHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken := r.FormValue("commenterToken")
	domain := r.Header.Get("Referer")

	if commenterToken == "" {
		fmt.Fprintf(w, "Error: %s\n", errorMissingField.Error())
		return
	}

	domain = domainStrip(domain)
	if domain == "" {
		fmt.Fprintf(w, "Error: No Referer header found in request\n")
		return
	}

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != errorNoSuchToken {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	d, err := domainGet(domain)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", errorNoSuchDomain.Error())
		return
	}

	if !d.SsoProvider {
		fmt.Fprintf(w, "Error: SSO not configured for %s\n", domain)
		return
	}

	if d.SsoSecret == "" || d.SsoUrl == "" {
		fmt.Fprintf(w, "Error: %s\n", errorMissingConfig.Error())
		return
	}

	key, err := hex.DecodeString(d.SsoSecret)
	if err != nil {
		logger.Errorf("cannot decode SSO secret as hex: %v", err)
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	token, err := ssoTokenNew(domain, commenterToken)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		logger.Errorf("cannot decode hex token: %v", err)
		fmt.Fprintf(w, "Error: %s\n", errorInternal.Error())
		return
	}

	h := hmac.New(sha256.New, key)
	h.Write(tokenBytes)
	signature := hex.EncodeToString(h.Sum(nil))

	u, err := url.Parse(d.SsoUrl)
	if err != nil {
		// this should really not be happening; we're checking if the
		// passed URL is valid at domain update
		logger.Errorf("cannot parse URL: %v", err)
		fmt.Fprintf(w, "Error: %s\n", errorInternal.Error())
		return
	}

	q := u.Query()
	q.Set("token", token)
	q.Set("hmac", signature)
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusFound)
}
