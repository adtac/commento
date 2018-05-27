package main

import (
	"testing"
)

func TestCommenterIsProviderUserBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google")

	exists, err := commenterIsProviderUser("google", "test@example.com")
	if err != nil {
		t.Errorf("unexpected error checking if commenter is a provider user: %v", err)
		return
	}

	if !exists {
		t.Errorf("user expected to exist not found")
		return
	}

	exists, err = commenterIsProviderUser("google", "test2@example.com")
	if err != nil {
		t.Errorf("unexpected error checking if commenter is a provider user: %v", err)
		return
	}

	if exists {
		t.Errorf("user expected to not exist not found")
		return
	}
}

func TestCommenterIsProviderUserEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterIsProviderUser("google", ""); err == nil {
		t.Errorf("expected error not found when checking for user with empty email")
		return
	}
}
