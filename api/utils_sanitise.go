package main

import (
	"regexp"
)

var prePlusMatch = regexp.MustCompile(`([^@\+]*)\+?(.*)@.*`)
var periodsMatch = regexp.MustCompile(`[\.]`)
var postAtMatch = regexp.MustCompile(`[^@]*(@.*)`)

func stripEmail(email string) string {
	postAt := postAtMatch.ReplaceAllString(email, `$1`)
	prePlus := prePlusMatch.ReplaceAllString(email, `$1`)
	strippedEmail := periodsMatch.ReplaceAllString(prePlus, ``) + postAt

	return strippedEmail
}

var https = regexp.MustCompile(`(https?://)`)
var trailingSlash = regexp.MustCompile(`(/*$)`)

func stripDomain(domain string) string {
	noSlash := trailingSlash.ReplaceAllString(domain, ``)
	noProtocol := https.ReplaceAllString(noSlash, ``)

	return noProtocol
}

var pathMatch = regexp.MustCompile(`(https?://[^/]*)`)

func stripPath(url string) string {
	strippedPath := pathMatch.ReplaceAllString(url, ``)

	return strippedPath
}

var httpsUrl = regexp.MustCompile(`^https?://`)

func isHttpsUrl(in string) bool {
	// Admittedly, this isn't the greatest URL checker. But it does what we need.
	// I don't care if the user gives an invalid URL, I just want to make sure
	// they don't do any XSS shenanigans. Hopefully, enforcing a https?:// prefix
	// solves this. If this function returns false, prefix with "http://"
	return len(httpsUrl.FindAllString(in, -1)) != 0
}
