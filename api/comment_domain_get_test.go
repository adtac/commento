package main

import (
	"testing"
	"time"
)

func TestCommentDomainGetBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commentHex, _ := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	domain, err := commentDomainGet(commentHex)
	if err != nil {
		t.Errorf("unexpected error getting domain by hex: %v", err)
		return
	}

	if domain != "example.com" {
		t.Errorf("expected domain = example.com got domain = %s", domain)
		return
	}
}

func TestCommentDomainGetEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commentDomainGet(""); err == nil {
		t.Errorf("expected error not found getting domain with empty commentHex")
		return
	}
}
