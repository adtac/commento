package main

import (
	"testing"
	"time"
)

func TestCommentOwnershipVerifyBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commentHex, _ := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	isOwner, err := commentOwnershipVerify("temp-commenter-hex", commentHex)
	if err != nil {
		t.Errorf("unexpected error verifying ownership: %v", err)
		return
	}

	if !isOwner {
		t.Errorf("expected to be owner of comment")
		return
	}

	isOwner, err = commentOwnershipVerify("another-commenter-hex", commentHex)
	if err != nil {
		t.Errorf("unexpected error verifying ownership: %v", err)
		return
	}

	if isOwner {
		t.Errorf("unexpected owner of comment not created by another-commenter-hex")
		return
	}
}

func TestCommentOwnershipVerifyEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commentOwnershipVerify("temp-commenter-hex", ""); err == nil {
		t.Errorf("expected error not founding verifying ownership with empty commentHex")
		return
	}
}
