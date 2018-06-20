package main

import (
	"testing"
)

func TestOwnerLoginBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := ownerLogin("test@example.com", "hunter2"); err == nil {
		t.Errorf("expected error not found when logging in without creating an account")
		return
	}

	ownerNew("test@example.com", "Test", "hunter2")

	if _, err := ownerLogin("test@example.com", "hunter2"); err != nil {
		t.Errorf("unexpected error when logging in: %v", err)
		return
	}

	if _, err := ownerLogin("test@example.com", "h******"); err == nil {
		t.Errorf("expected error not found when given wrong password")
		return
	}

	if ownerToken, err := ownerLogin("test@example.com", "hunter2"); ownerToken == "" {
		t.Errorf("empty token on successful login: %v", err)
		return
	}
}

func TestOwnerLoginEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := ownerLogin("test@example.com", ""); err == nil {
		t.Errorf("expected error not found when passing empty password")
		return
	}

	ownerNew("test@example.com", "Test", "hunter2")

	if _, err := ownerLogin("test@example.com", ""); err == nil {
		t.Errorf("expected error not found when passing empty password")
		return
	}
}
