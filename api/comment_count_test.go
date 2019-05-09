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

	counts, err := commentCount("example.com", []string{"/path.html"})
	if err != nil {
		t.Errorf("unexpected error counting comments: %v", err)
		return
	}

	if counts["/path.html"] != 3 {
		t.Errorf("expected count=3 got count=%d", counts["/path.html"])
		return
	}
}

func TestCommentCountNewPage(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	counts, err := commentCount("example.com", []string{"/path.html"})
	if err != nil {
		t.Errorf("unexpected error counting comments: %v", err)
		return
	}

	if counts["/path.html"] != 0 {
		t.Errorf("expected count=0 got count=%d", counts["/path.html"])
		return
	}
}

func TestCommentCountEmpty(t *testing.T) {
	if _, err := commentCount("example.com", []string{""}); err != nil {
		t.Errorf("unexpected error counting comments on empty path: %v", err)
		return
	}

	if _, err := commentCount("", []string{""}); err == nil {
		t.Errorf("expected error not found counting comments with empty everything")
		return
	}
}
