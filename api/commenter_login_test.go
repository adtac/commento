package main

import (
	"testing"
)

func TestCommenterLoginBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterLogin("test@example.com", "hunter2"); err == nil {
		t.Errorf("expected error not found when logging in without creating an account")
		return
	}

	commenterNew("test@example.com", "Test", "undefined", "undefined", "commento", "hunter2")

	if _, err := commenterLogin("test@example.com", "hunter2"); err != nil {
		t.Errorf("unexpected error when logging in: %v", err)
		return
	}

	if _, err := commenterLogin("test@example.com", "h******"); err == nil {
		t.Errorf("expected error not found when given wrong password")
		return
	}

	if commenterToken, err := commenterLogin("test@example.com", "hunter2"); commenterToken == "" {
		t.Errorf("empty comenterToken on successful login: %v", err)
		return
	}
}

func TestCommenterLoginEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterLogin("test@example.com", ""); err == nil {
		t.Errorf("expected error not found when passing empty password")
		return
	}

	commenterNew("test@example.com", "Test", "undefined", "", "commenter", "hunter2")

	if _, err := commenterLogin("test@example.com", ""); err == nil {
		t.Errorf("expected error not found when passing empty password")
		return
	}
}

func TestCommenterLoginNonCommento(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterNew("test@example.com", "Test", "undefined", "undefined", "google", "")

	if _, err := commenterLogin("test@example.com", "hunter2"); err == nil {
		t.Errorf("expected error not found logging into a non-Commento account")
		return
	}
}
