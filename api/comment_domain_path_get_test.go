package main

import (
	"testing"
	"time"
)

func TestCommentDomainPathGetBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commentHex, _ := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	domain, path, err := commentDomainPathGet(commentHex)
	if err != nil {
		t.Errorf("unexpected error getting domain by hex: %v", err)
		return
	}

	if domain != "example.com" {
		t.Errorf("expected domain=example.com got domain=%s", domain)
		return
	}

	if path != "/path.html" {
		t.Errorf("expected path=/path.html got path=%s", path)
		return
	}
}

func TestCommentDomainGetEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, _, err := commentDomainPathGet(""); err == nil {
		t.Errorf("expected error not found getting domain with empty commentHex")
		return
	}
}
