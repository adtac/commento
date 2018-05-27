package main

import (
	"testing"
)

func TestCommenterSessionUpdateBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	session, _ := commenterSessionNew()

	if err := commenterSessionUpdate(session, "temp-commenter-hex"); err != nil {
		t.Errorf("unexpected error updating session to commenterHex: %v", err)
		return
	}
}

func TestCommenterSessionUpdateEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if err := commenterSessionUpdate("", "temp-commenter-hex"); err == nil {
		t.Errorf("expected error not found when updating with empty session")
		return
	}
}
