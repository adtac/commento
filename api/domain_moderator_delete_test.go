package main

import (
	"testing"
)

func TestDomainModeratorDeleteBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	domainModeratorNew("example.com", "test@example.com")
	domainModeratorNew("example.com", "test2@example.com")

	if err := domainModeratorDelete("example.com", "test@example.com"); err != nil {
		t.Errorf("unexpected error creating new domain moderator: %v", err)
		return
	}

	isMod, _ := isDomainModerator("example.com", "test@example.com")
	if isMod {
		t.Errorf("email %s still moderator after deletion", "test@example.com")
		return
	}

	isMod, _ = isDomainModerator("example.com", "test2@example.com")
	if !isMod {
		t.Errorf("email %s no longer moderator after deleting a different email", "test@example.com")
		return
	}
}

func TestDomainModeratorDeleteEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	domainModeratorNew("example.com", "test@example.com")

	if err := domainModeratorDelete("example.com", ""); err == nil {
		t.Errorf("expected error not found when passing empty email")
		return
	}

	if err := domainModeratorDelete("", ""); err == nil {
		t.Errorf("expected error not found when passing empty everything")
		return
	}
}
