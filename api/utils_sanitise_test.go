package main

import (
	"testing"
)

func TestEmailStripBasics(t *testing.T) {
	tests := map[string]string{
		"test@example.com":              "test@example.com",
		"test+strip@example.com":        "test@example.com",
		"test+strip+strip2@example.com": "test@example.com",
	}

	for in, out := range tests {
		if emailStrip(in) != out {
			t.Errorf("for in=%s expected out=%s got out=%s", in, out, emailStrip(in))
			return
		}
	}
}
