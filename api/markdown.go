package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var policy *bluemonday.Policy
var renderer blackfriday.Renderer
var extensions int

func markdownRendererCreate() error {
	policy = bluemonday.UGCPolicy()

	extensions = 0
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH

	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_SKIP_HTML
	htmlFlags |= blackfriday.HTML_SKIP_IMAGES
	htmlFlags |= blackfriday.HTML_SAFELINK
	htmlFlags |= blackfriday.HTML_HREF_TARGET_BLANK

	renderer = blackfriday.HtmlRenderer(htmlFlags, "", "")

	return nil
}
