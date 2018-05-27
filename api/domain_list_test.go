package main

import (
	"testing"
)

func TestDomainListBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	domainNew("temp-owner-hex", "Example", "example.com")
	domainNew("temp-owner-hex", "Example", "example2.com")

	d, err := domainList("temp-owner-hex")
	if err != nil {
		t.Errorf("unexpected error listing domains: %v", err)
		return
	}

	if len(d) != 2 {
		t.Errorf("expected number of domains to be 2 got %d", len(d))
		return
	}

	if d[0].Domain != "example.com" {
		t.Errorf("expected first domain to be example.com got %s", d[0].Domain)
		return
	}

	if d[1].Domain != "example2.com" {
		t.Errorf("expected first domain to be example2.com got %s", d[1].Domain)
		return
	}
}
