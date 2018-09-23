package main

import (
	"testing"
	"time"
)

func TestCommentCountBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "http://example.com/photo.jpg", "google", "")

	commentNew(commenterHex, "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())
	commentNew(commenterHex, "example.com", "/path.html", "root", "**bar**", "approved", time.Now().UTC())
	commentNew(commenterHex, "example.com", "/path.html", "root", "**baz**", "unapproved", time.Now().UTC())

	count, err := commentCount("example.com", "/path.html")
	if err != nil {
		t.Errorf("unexpected error counting comments: %v", err)
		return
	}

	if count != 2 {
		t.Errorf("expected count=2 got count=%d", count)
		return
	}
}

func TestCommentCountNewPage(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	count, err := commentCount("example.com", "/path.html")
	if err != nil {
		t.Errorf("unexpected error counting comments: %v", err)
		return
	}

	if count != 0 {
		t.Errorf("expected count=0 got count=%d", count)
		return
	}
}

func TestCommentCountEmpty(t *testing.T) {
	if _, err := commentCount("example.com", ""); err != nil {
		t.Errorf("unexpected error counting comments on empty path: %v", err)
		return
	}

	if _, err := commentCount("", ""); err == nil {
		t.Errorf("expected error not found counting comments with empty everything")
		return
	}
}
