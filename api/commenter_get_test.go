package main

import (
	"testing"
)

func TestCommenterGetByHexBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")

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

func TestCommenterGetByCommenterToken(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")

	commenterToken, _ := commenterTokenNew()

	commenterSessionUpdate(commenterToken, commenterHex)

	c, err := commenterGetByCommenterToken(commenterToken)
	if err != nil {
		t.Errorf("unexpected error getting commenter by hex: %v", err)
		return
	}

	if c.Name != "Test" {
		t.Errorf("expected name=Test got name=%s", c.Name)
		return
	}
}

func TestCommenterGetByCommenterTokenEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterGetByCommenterToken(""); err == nil {
		t.Errorf("expected error not found getting commenter with empty commenterToken")
		return
	}
}

func TestCommenterGetByName(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")

	commenterToken, _ := commenterTokenNew()

	commenterSessionUpdate(commenterToken, commenterHex)

	c, err := commenterGetByEmail("google", "test@example.com")
	if err != nil {
		t.Errorf("unexpected error getting commenter by hex: %v", err)
		return
	}

	if c.Name != "Test" {
		t.Errorf("expected name=Test got name=%s", c.Name)
		return
	}
}

func TestCommenterGetByNameEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterGetByEmail("", ""); err == nil {
		t.Errorf("expected error not found getting commenter with empty everything")
		return
	}
}
