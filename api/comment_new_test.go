package main

import (
	"testing"
	"time"
)

func TestCommentNewBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC()); err != nil {
		t.Errorf("unexpected error creating new comment: %v", err)
		return
	}
}

func TestCommentNewEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := commentNew("temp-commenter-hex", "example.com", "", "root", "**foo**", "approved", time.Now().UTC()); err != nil {
		t.Errorf("empty path not allowed: %v", err)
		return
	}

	if _, err := commentNew("temp-commenter-hex", "", "", "root", "**foo**", "approved", time.Now().UTC()); err == nil {
		t.Errorf("expected error not found creatingn new comment with empty domain")
		return
	}

	if _, err := commentNew("", "", "", "", "", "", time.Now().UTC()); err == nil {
		t.Errorf("expected error not found creatingn new comment with empty everything")
		return
	}
}

func TestCommentNewUpvoted(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commentHex, _ := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	statement := `
		SELECT score
		FROM comments
		WHERE commentHex = $1;
	`
	row := db.QueryRow(statement, commentHex)

	var score int
	if err := row.Scan(&score); err != nil {
		t.Errorf("error scanning score from comments table: %v", err)
		return
	}

	if score != 0 {
		t.Errorf("expected comment to be at 0 points")
		return
	}
}

func TestCommentNewThreadLocked(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	pageNew("example.com", "/path.html")
	p, _ := pageGet("example.com", "/path.html")
	p.IsLocked = true
	pageUpdate(p)

	_, err := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())
	if err == nil {
		t.Errorf("expected error not found creating a new comment on a locked thread")
		return
	}
}
