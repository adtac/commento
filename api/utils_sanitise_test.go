package main

import (
	"testing"
)

func TestStripEmailBasics(t *testing.T) {
	tests := map[string]string{
		"test@example.com":              "test@example.com",
		"test+strip@example.com":        "test@example.com",
		"test+strip+strip2@example.com": "test@example.com",
	}

	for in, out := range tests {
		if stripEmail(in) != out {
			t.Errorf("for in=%s expected out=%s got out=%s", in, out, stripEmail(in))
			return
		}
	}
}
