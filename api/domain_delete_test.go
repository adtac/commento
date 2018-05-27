package main

import (
	"testing"
)

func TestDomainDeleteBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	domainNew("temp-owner-hex", "Example", "example.com")
	domainNew("temp-owner-hex", "Example", "example2.com")

	if err := domainDelete("example.com"); err != nil {
		t.Errorf("unexpected error deleting domain: %v", err)
		return
	}

	d, _ := domainList("temp-owner-hex")

	if len(d) != 1 {
		t.Errorf("expected number of domains to be 1 got %d", len(d))
		return
	}

	if d[0].Domain != "example2.com" {
		t.Errorf("expected first domain to be example2.com got %s", d[0].Domain)
		return
	}
}

func TestDomainDeleteEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if err := domainDelete(""); err == nil {
		t.Errorf("expected error not found when deleting with empty domain")
		return
	}
}
