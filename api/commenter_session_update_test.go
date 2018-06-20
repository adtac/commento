package main

import (
	"testing"
)

func TestCommenterSessionUpdateBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterToken, _ := commenterTokenNew()

	if err := commenterSessionUpdate(commenterToken, "temp-commenter-hex"); err != nil {
		t.Errorf("unexpected error updating commenter session: %v", err)
		return
	}
}

func TestCommenterSessionUpdateEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if err := commenterSessionUpdate("", "temp-commenter-hex"); err == nil {
		t.Errorf("expected error not found when updating with empty commenterToken")
		return
	}
}
