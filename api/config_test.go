package main

import (
	"os"
	"testing"
)

func TestParseConfigBasics(t *testing.T) {
	os.Setenv("COMMENTO_ORIGIN", "https://commento.io")

	if err := parseConfig(); err != nil {
		t.Errorf("unexpected error when parsing config: %v", err)
		return
	}

	// This test feels kinda stupid, but whatever.
	if os.Getenv("PORT") != "8080" {
		t.Errorf("expected PORT=8080, but PORT=%s instead", os.Getenv("PORT"))
		return
	}

	os.Setenv("COMMENTO_PORT", "1886")

	if err := parseConfig(); err != nil {
		t.Errorf("unexpected error when parsing config: %v", err)
		return
	}

	if os.Getenv("PORT") != "1886" {
		t.Errorf("expected PORT=1886, but PORT=%s instead", os.Getenv("PORT"))
		return
	}
}
