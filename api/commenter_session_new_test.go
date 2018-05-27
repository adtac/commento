package main

import (
	"testing"
)

func TestCommenterSessionNewBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterSessionNew(); err != nil {
		t.Errorf("unexpected error creating new session: %v", err)
		return
	}
}
