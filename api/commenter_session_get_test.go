package main

import (
	"testing"
)

func TestCommenterSessionGetBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google")

	session, _ := commenterSessionNew()

	commenterSessionUpdate(session, commenterHex)

	cs, err := commenterSessionGet(session)
	if err != nil {
		t.Errorf("unexpected error found when getting session information: %v", err)
		return
	}

	if cs.CommenterHex != commenterHex {
		t.Errorf("expected commenterHex=%s got commenterHex=%s", commenterHex, cs.CommenterHex)
		return
	}
}

func TestCommenterSessionGetDNE(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	_, err := commenterSessionGet("does-not-exist")
	if err == nil {
		t.Errorf("expected error not found when invalid session")
		return
	}
}

func TestCommenterSessionGetEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	_, err := commenterSessionGet("")
	if err == nil {
		t.Errorf("expected error not found with empty session")
		return
	}
}
