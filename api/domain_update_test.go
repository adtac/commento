package main

import (
	"testing"
)

func TestDomainUpdateBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	domainNew("temp-owner-hex", "Example", "example.com")

	d, _ := domainList("temp-owner-hex")

	d[0].Name = "Example2"

	if err := domainUpdate(d[0]); err != nil {
		t.Errorf("unexpected error updating domain: %v", err)
		return
	}

	d, _ = domainList("temp-owner-hex")

	if d[0].Name != "Example2" {
		t.Errorf("expected name=Example2 got name=%s", d[0].Name)
		return
	}
}
