package main

import (
	"testing"
)

func TestPageGetBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	pageNew("example.com", "/path.html")

	p, err := pageGet("example.com", "/path.html")
	if err != nil {
		t.Errorf("unexpected error getting page: %v", err)
		return
	}

	if p.IsLocked != false {
		t.Errorf("expected p.IsLocked=false got %v", p.IsLocked)
		return
	}

	if _, err := pageGet("example.com", "/path2.html"); err != nil {
		t.Errorf("unexpected error getting page with non-existant record: %v", err)
		return
	}
}

func TestPageGetEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	pageNew("example.com", "")

	if _, err := pageGet("example.com", ""); err != nil {
		t.Errorf("unexpected error getting page with empty path: %v", err)
		return
	}

	if _, err := pageGet("", "/path.html"); err == nil {
		t.Errorf("exepected error not found when getting page with empty domain")
		return
	}
}
