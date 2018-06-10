package main

import (
	"testing"
	"time"
)

func TestCommentVoteBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	cr0, _ := commenterNew("test1@example.com", "Test1", "undefined", "http://example.com/photo.jpg", "google", "")
	cr1, _ := commenterNew("test2@example.com", "Test2", "undefined", "http://example.com/photo.jpg", "google", "")
	cr2, _ := commenterNew("test3@example.com", "Test3", "undefined", "http://example.com/photo.jpg", "google", "")

	c0, _ := commentNew(cr0, "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	commentVote(cr0, c0, -1)
	if c, _, _ := commentList("temp", "example.com", "/path.html", false); c[0].Score != -1 {
		t.Errorf("expected c[0].Score = -1 got c[0].Score = %d", c[0].Score)
		return
	}

	commentVote(cr1, c0, -1)
	commentVote(cr2, c0, -1)
	if c, _, _ := commentList("temp", "example.com", "/path.html", false); c[0].Score != -3 {
		t.Errorf("expected c[0].Score = -3 got c[0].Score = %d", c[0].Score)
		return
	}

	commentVote(cr1, c0, -1)
	if c, _, _ := commentList("temp", "example.com", "/path.html", false); c[0].Score != -3 {
		t.Errorf("expected c[0].Score = -3 got c[0].Score = %d", c[0].Score)
		return
	}

	commentVote(cr1, c0, 0)
	if c, _, _ := commentList("temp", "example.com", "/path.html", false); c[0].Score != -2 {
		t.Errorf("expected c[0].Score = -2 got c[0].Score = %d", c[0].Score)
		return
	}

	c1, _ := commentNew(cr1, "example.com", "/path.html", "root", "**bar**", "approved", time.Now().UTC())

	commentVote(cr0, c1, 0)
	if c, _, _ := commentList("temp", "example.com", "/path.html", false); c[1].Score != 1 {
		t.Errorf("expected c[1].Score = 1 got c[1].Score = %d", c[1].Score)
		return
	}

	commentVote(cr1, c1, 0)
	if c, _, _ := commentList("temp", "example.com", "/path.html", false); c[1].Score != 0 {
		t.Errorf("expected c[1].Score = 0 got c[1].Score = %d", c[1].Score)
		return
	}
}
