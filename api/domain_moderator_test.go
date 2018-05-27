package main

import (
	"testing"
)

func TestDomainModeratorListBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	domainModeratorNew("example.com", "test@example.com")
	domainModeratorNew("example.com", "test2@example.com")

	mods, err := domainModeratorList("example.com")
	if err != nil {
		t.Errorf("unexpected error listing domain moderators: %v", err)
		return
	}

	if len(mods) != 2 {
		t.Errorf("expected number of domain moderators to be 2 got %d", len(mods))
		return
	}

	if mods[0].Email != "test@example.com" {
		t.Errorf("expected first domain to be test@example.com got %s", mods[0].Email)
		return
	}

	if mods[1].Email != "test2@example.com" {
		t.Errorf("expected first domain to be test2@example.com got %s", mods[0].Email)
		return
	}
}

func TestIsDomainModeratorBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	domainModeratorNew("example.com", "test@example.com")

	isMod, err := isDomainModerator("example.com", "test@example.com")
	if err != nil {
		t.Errorf("unexpected error checking if email is a moderator: %v", err)
		return
	}

	if !isMod {
		t.Errorf("expected test@example.com to be a moderator got isMod=false")
		return
	}

	isMod, err = isDomainModerator("example.com", "test2@example.com")
	if err != nil {
		t.Errorf("unexpected error checking if email is a moderator: %v", err)
		return
	}

	if isMod {
		t.Errorf("expected test2@example.com to not be a moderator got isMod=true")
		return
	}
}
