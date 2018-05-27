package main

import (
	"testing"
)

func TestDomainGetBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	domainNew("temp-owner-hex", "Example", "example.com")

	d, err := domainGet("example.com")
	if err != nil {
		t.Errorf("unexpected error getting domain: %v", err)
		return
	}

	if d.Name != "Example" {
		t.Errorf("expected name=Example got name=%s", d.Name)
		return
	}
}

func TestDomainGetEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := domainGet(""); err == nil {
		t.Errorf("expected error not found when getting with empty domain")
		return
	}
}

func TestDomainGetDNE(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := domainGet("example.com"); err == nil {
		t.Errorf("expected error not found when getting non-existant domain")
		return
	}
}
