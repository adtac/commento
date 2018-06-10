package main

import (
	"testing"
)

func TestCommenterNewBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", ""); err != nil {
		t.Errorf("unexpected error creating new commenter: %v", err)
		return
	}
}

func TestCommenterNewEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterNew("", "Test", "undefined", "https://example.com/photo.jpg", "google", ""); err == nil {
		t.Errorf("expected error not found creating new commenter with empty email")
		return
	}

	if _, err := commenterNew("", "", "", "", "", ""); err == nil {
		t.Errorf("expected error not found creating new commenter with empty everything")
		return
	}
}

func TestCommenterNewCommento(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterNew("test@example.com", "Test", "undefined", "", "commento", ""); err == nil {
		t.Errorf("expected error not found creating new commento account with empty password")
		return
	}
}
