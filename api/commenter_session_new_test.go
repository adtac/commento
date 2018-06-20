package main

import (
	"testing"
)

func TestCommenterTokenNewBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commenterTokenNew(); err != nil {
		t.Errorf("unexpected error creating new commenterToken: %v", err)
		return
	}
}
