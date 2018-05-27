package main

import (
	"testing"
)

func TestOwnerSelfBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	ownerNew("test@example.com", "Test", "hunter2")
	session, _ := ownerLogin("test@example.com", "hunter2")

	loggedIn, o := ownerSelf(session)
	if !loggedIn {
		t.Errorf("expected loggedIn=true got loggedIn=false")
		return
	}

	if o.Name != "Test" {
		t.Errorf("expected name=Test got name=%s", o.Name)
		return
	}
}

func TestOwnerSelfNotLoggedIn(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if loggedIn, _ := ownerSelf("does-not-exist"); loggedIn {
		t.Errorf("expected loggedIn=false got loggedIn=true")
		return
	}
}
