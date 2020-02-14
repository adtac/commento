package main

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

func commenterPhotoHandler(w http.ResponseWriter, r *http.Request) {
	c, err := commenterGetByHex(r.FormValue("commenterHex"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	url := c.Photo
	if c.Provider == "google" {
		if strings.HasSuffix(url, "photo.jpg") {
			url += "?sz=38"
		} else {
			url += "=s38"
		}
	} else if c.Provider == "github" {
		url += "&s=38"
	} else if c.Provider == "twitter" {
		url += "?size=normal"
	} else if c.Provider == "gitlab" {
		url += "?width=38"
	}

	resp, err := http.Get(url)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer resp.Body.Close()

	if c.Provider != "commento" { // Custom URL avatars need to be resized.
		io.Copy(w, resp.Body)
		return
	}

	// Limit the size of the response to 128 KiB to prevent DoS attacks
	// that exhaust memory.
	limitedResp := &io.LimitedReader{R: resp.Body, N: 128 * 1024}

	img, err := jpeg.Decode(limitedResp)
	if err != nil {
		fmt.Fprintf(w, "JPEG decode failed: %v\n", err)
		return
	}

	if err = imaging.Encode(w, imaging.Resize(img, 38, 0, imaging.Lanczos), imaging.JPEG); err != nil {
		fmt.Fprintf(w, "image encoding failed: %v\n", err)
		return
	}
}
