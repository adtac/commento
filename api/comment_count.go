package main

import (
	"net/http"
)

func commentCount(domain string, path string) (int, error) {
	// path can be empty
	if domain == "" {
		return 0, errorMissingField
	}

	p, err := pageGet(domain, path)
	if err != nil {
		return 0, errorInternal
	}

	return p.CommentCount, nil
}

func commentCountHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Domain         *string `json:"domain"`
		Path           *string `json:"path"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	path := *x.Path

	count, err := commentCount(domain, path)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "count": count})
}
