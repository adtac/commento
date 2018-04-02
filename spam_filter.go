package main

import (
	"github.com/adtac/go-akismet/akismet"
	"net/http"
	"os"
	"regexp"
)

func checkSpam(r *http.Request, url string, name string, comment string) bool {
	akismetKey := os.Getenv("AKISMET_KEY")
	if akismetKey == "" {
		return false
	}

	IP := r.RemoteAddr
	if r.Header.Get("X-Forwarded-For") != "" {
		IP = r.Header.Get("X-Forwarded-For")
	}

	exp := regexp.MustCompile(`^(?P<Blog>(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(?:[^:\/\n]+))`)
	match := exp.FindStringSubmatch(url)

	var blog string
	for i, name := range exp.SubexpNames() {
		if name == "Blog" {
			blog = match[i]
		}
	}

	isSpam, err := akismet.Check(&akismet.Comment{
		Blog:           blog,
		UserIP:         IP,
		UserAgent:      r.Header.Get("User-Agent"),
		CommentType:    "comment",
		CommentAuthor:  name,
		CommentContent: comment,
	}, akismetKey)

	if err != nil {
		Logger.Errorf("error: cannot validate commenet using Akismet: %v", err)
		return true
	}

	return isSpam
}
