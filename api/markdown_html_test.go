package main

import (
	"strings"
	"testing"
)

func TestMarkdownToHtmlBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	// basic markdown and expected html tests
	tests := map[string]string{
		"Foo": "<p>Foo</p>",

		"Foo\n\nBar": "<p>Foo</p>\n\n<p>Bar</p>",

		"XSS: <script src='http://example.com/script.js'></script> Foo": "<p>XSS:  Foo</p>",

		"Regular [Link](http://example.com)": "<p>Regular <a href=\"http://example.com\" rel=\"nofollow\">Link</a></p>",

		"XSS [Link](data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pgo=)": "<p>XSS <tt>Link</tt></p>",

		"![Images disallowed](http://example.com/image.jpg)": "<p></p>",

		"**bold** *italics*": "<p><strong>bold</strong> <em>italics</em></p>",

		"http://example.com/autolink": "<p><a href=\"http://example.com/autolink\" rel=\"nofollow\">http://example.com/autolink</a></p>",

		"<b>not bold</b>": "<p>not bold</p>",
	}

	for in, out := range tests {
		html := strings.TrimSpace(markdownToHtml(in))
		if html != out {
			t.Errorf("for in=[%s] expected out=[%s] got out=[%s]", in, out, html)
			return
		}
	}
}
