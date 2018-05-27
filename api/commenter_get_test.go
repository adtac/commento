package main

import (
	"testing"
)

func TestCommenterGetByHexBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google")

	c, err := commenterGetByHex(commenterHex)
	if err != nil {
		t.Errorf("unexpected error getting commenter by hex: %v", err)
		return
	}

	if c.Name != "Test" {
		t.Errorf("expected name=Test got name=%s", c.Name)
		return
	}
}

func TestCommenterGetByHexEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterGetByHex(""); err == nil {
		t.Errorf("expected error not found getting commenter with empty hex")
		return
	}
}

func TestCommenterGetBySession(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google")

	session, _ := commenterSessionNew()

	commenterSessionUpdate(session, commenterHex)

	c, err := commenterGetBySession(session)
	if err != nil {
		t.Errorf("unexpected error getting commenter by hex: %v", err)
		return
	}

	if c.Name != "Test" {
		t.Errorf("expected name=Test got name=%s", c.Name)
		return
	}
}

func TestCommenterGetBySessionEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterGetBySession(""); err == nil {
		t.Errorf("expected error not found getting commenter with empty session")
		return
	}
}
