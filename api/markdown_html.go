package main

import (
	"gopkg.in/russross/blackfriday.v1"
)

func markdownToHtml(markdown string) string {
	unsafe := blackfriday.Markdown([]byte(markdown), renderer, extensions)
	return string(policy.SanitizeBytes(unsafe))
}
