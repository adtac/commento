package main

import (
	"os"
	"testing"
	"path/filepath"
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

func TestParseConfigStatic(t *testing.T) {
	os.Setenv("COMMENTO_ORIGIN", "https://commento.io")

	if err := parseConfig(); err != nil {
		t.Errorf("unexpected error when parsing config: %v", err)
		return
	}

	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Errorf("cannot load binary path: %v", err)
		return
	}

	if os.Getenv("STATIC") != binPath {
		t.Errorf("COMMENTO_STATIC != %s when unset", binPath)
		return
	}

	os.Setenv("COMMENTO_STATIC", "/usr/")

	if err := parseConfig(); err != nil {
		t.Errorf("unexpected error when parsing config: %v", err)
		return
	}

	if os.Getenv("STATIC") != "/usr" {
		t.Errorf("COMMENTO_STATIC != /usr when unset")
		return
	}
}

func TestParseConfigStaticDNE(t *testing.T) {
	os.Setenv("COMMENTO_ORIGIN", "https://commento.io")
	os.Setenv("COMMENTO_STATIC", "/does/not/exist/surely/")

	if err := parseConfig(); err == nil {
		t.Errorf("expected error not found when a non-existant directory is used")
		return
	}
}

func TestParseConfigStaticNotADirectory(t *testing.T) {
	os.Setenv("COMMENTO_ORIGIN", "https://commento.io")
	os.Setenv("COMMENTO_STATIC", os.Args[0])

	if err := parseConfig(); err != errorNotADirectory {
		t.Errorf("expected error not found when a file is used")
		return
	}
}
