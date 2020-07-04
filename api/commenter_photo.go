package main

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"strings"
	"strconv"

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

	// Limit the size of the response to 128 KiB or configured size
	// to prevent DoS attacks that exhaust memory.
	var limit int
	limitString := os.Getenv("IMAGE_SIZE_LIMIT")
	if (limitString == "") {
		limit = 128 * 1024
	} else {
		limit = strconv.Atoi(limitString)
	}
	limitedResp := &io.LimitedReader{R: resp.Body, N: limit}

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
