package main

import (
	"golang.org/x/net/html"
	"net/http"
)

func htmlTitleRecurse(h *html.Node) string {
	if h == nil || h.FirstChild == nil {
		return ""
	}

	if h.Type == html.ElementNode && h.Data == "title" {
		return h.FirstChild.Data
	}

	for c := h.FirstChild; c != nil; c = c.NextSibling {
		res := htmlTitleRecurse(c)
		if res != "" {
			return res
		}
	}

	return ""
}

func htmlTitleGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	h, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	return htmlTitleRecurse(h), nil
}
