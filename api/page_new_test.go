package main

import (
	"testing"
)

func TestPageNewBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if err := pageNew("example.com", "/path.html"); err != nil {
		t.Errorf("unexpected error creating page: %v", err)
		return
	}
}

func TestPageNewEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if err := pageNew("example.com", ""); err != nil {
		t.Errorf("unexpected error creating page with empty path: %v", err)
		return
	}

	if err := pageNew("", "/path.html"); err == nil {
		t.Errorf("expected error not found creating page with empty domain")
		return
	}
}

func TestPageNewUnique(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if err := pageNew("example.com", "/path.html"); err != nil {
		t.Errorf("unexpected error creating page: %v", err)
		return
	}

	// no error should be returned when trying to duplicate insert
	if err := pageNew("example.com", "/path.html"); err != nil {
		t.Errorf("unexpected error creating same page twice: %v", err)
		return
	}
}
