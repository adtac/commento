package main

import (
	"strings"
	"testing"
	"time"
)

func TestCommentListBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "http://example.com/photo.jpg", "google", "")

	commentNew(commenterHex, "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())
	commentNew(commenterHex, "example.com", "/path.html", "root", "**bar**", "approved", time.Now().UTC())

	c, _, err := commentList("temp-commenter-hex", "example.com", "/path.html", false)
	if err != nil {
		t.Errorf("unexpected error listing page comments: %v", err)
		return
	}

	if len(c) != 2 {
		t.Errorf("expected 2 comments got %d comments", len(c))
		return
	}

	if c[0].Direction != 0 {
		t.Errorf("expected c.Direction = 0 got c.Direction = %d", c[0].Direction)
		return
	}

	c1Html := strings.TrimSpace(c[1].Html)
	if c1Html != "<p><strong>bar</strong></p>" {
		t.Errorf("expected c[1].Html=[<p><strong>bar</strong></p>] got c[1].Html=[%s]", c1Html)
		return
	}

	c, _, err = commentList(commenterHex, "example.com", "/path.html", false)
	if err != nil {
		t.Errorf("unexpected error listing page comments: %v", err)
		return
	}

	if len(c) != 2 {
		t.Errorf("expected 2 comments got %d comments", len(c))
		return
	}

	if c[0].Direction != 0 {
		t.Errorf("expected c.Direction = 1 got c.Direction = %d", c[0].Direction)
		return
	}
}

func TestCommentListEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, _, err := commentList("temp-commenter-hex", "", "/path.html", false); err == nil {
		t.Errorf("expected error not found listing comments with empty domain")
		return
	}
}

func TestCommentListSelfUnapproved(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "http://example.com/photo.jpg", "google", "")

	commentNew(commenterHex, "example.com", "/path.html", "root", "**foo**", "unapproved", time.Now().UTC())

	c, _, _ := commentList("temp-commenter-hex", "example.com", "/path.html", false)

	if len(c) != 0 {
		t.Errorf("expected user to not see unapproved comment")
		return
	}

	c, _, _ = commentList(commenterHex, "example.com", "/path.html", false)

	if len(c) != 1 {
		t.Errorf("expected user to see unapproved self comment")
		return
	}
}

func TestCommentListAnonymousUnapproved(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commentNew("anonymous", "example.com", "/path.html", "root", "**foo**", "unapproved", time.Now().UTC())

	c, _, _ := commentList("anonymous", "example.com", "/path.html", false)

	if len(c) != 0 {
		t.Errorf("expected user to not see unapproved anonymous comment as anonymous")
		return
	}
}

func TestCommentListIncludeUnapproved(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commentNew("anonymous", "example.com", "/path.html", "root", "**foo**", "unapproved", time.Now().UTC())

	c, _, _ := commentList("anonymous", "example.com", "/path.html", true)

	if len(c) != 1 {
		t.Errorf("expected to see unapproved comments because includeUnapproved was true")
		return
	}
}

func TestCommentListDifferentPaths(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commentNew("anonymous", "example.com", "/path1.html", "root", "**foo**", "unapproved", time.Now().UTC())
	commentNew("anonymous", "example.com", "/path1.html", "root", "**foo**", "unapproved", time.Now().UTC())
	commentNew("anonymous", "example.com", "/path2.html", "root", "**foo**", "unapproved", time.Now().UTC())

	c, _, _ := commentList("anonymous", "example.com", "/path1.html", true)

	if len(c) != 2 {
		t.Errorf("expected len(c) = 2 got len(c) = %d", len(c))
		return
	}

	c, _, _ = commentList("anonymous", "example.com", "/path2.html", true)

	if len(c) != 1 {
		t.Errorf("expected len(c) = 1 got len(c) = %d", len(c))
		return
	}
}

func TestCommentListDifferentDomains(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commentNew("anonymous", "example1.com", "/path.html", "root", "**foo**", "unapproved", time.Now().UTC())
	commentNew("anonymous", "example2.com", "/path.html", "root", "**foo**", "unapproved", time.Now().UTC())

	c, _, _ := commentList("anonymous", "example1.com", "/path.html", true)

	if len(c) != 1 {
		t.Errorf("expected len(c) = 1 got len(c) = %d", len(c))
		return
	}

	c, _, _ = commentList("anonymous", "example2.com", "/path.html", true)

	if len(c) != 1 {
		t.Errorf("expected len(c) = 1 got len(c) = %d", len(c))
		return
	}
}
