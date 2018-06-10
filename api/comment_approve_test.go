package main

import (
	"testing"
	"time"
)

func TestCommentApproveBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")

	commentHex, _ := commentNew(commenterHex, "example.com", "/path.html", "root", "**foo**", "unapproved", time.Now().UTC())

	if err := commentApprove(commentHex); err != nil {
		t.Errorf("unexpected error approving comment: %v", err)
		return
	}

	if c, _, _ := commentList("anonymous", "example.com", "/path.html", true); c[0].State != "approved" {
		t.Errorf("expected state = approved got state = %s", c[0].State)
		return
	}
}

func TestCommentApproveEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if err := commentApprove(""); err == nil {
		t.Errorf("expected error not found approving comment with empty commentHex")
		return
	}
}
