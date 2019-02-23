package main

import (
	"io"
	"net/http"
)

func commenterPhotoHandler(w http.ResponseWriter, r *http.Request) {
	c, err := commenterGetByHex(r.FormValue("commenterHex"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	url := c.Photo
	if c.Provider == "google" {
		url += "?sz=50"
	} else if c.Provider == "github" {
		url += "&s=50"
	} else if c.Provider == "twitter" {
		url += "?size=normal"
	} else if c.Provider == "gitlab" {
		url += "?width=50"
	}

	resp, err := http.Get(url)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}
