package main

import (
	"testing"
)

func TestOwnerGetByEmailBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	ownerHex, _ := ownerNew("test@example.com", "Test", "hunter2")

	o, err := ownerGetByEmail("test@example.com")
	if err != nil {
		t.Errorf("unexpected error on ownerGetByEmail: %v", err)
		return
	}

	if o.OwnerHex != ownerHex {
		t.Errorf("expected ownerHex=%s got ownerHex=%s", ownerHex, o.OwnerHex)
		return
	}
}

func TestOwnerGetByEmailDNE(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := ownerGetByEmail("invalid@example.com"); err == nil {
		t.Errorf("expected error not found on ownerGetByEmail before creating an account")
		return
	}
}

func TestOwnerGetByOwnerTokenBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	ownerHex, _ := ownerNew("test@example.com", "Test", "hunter2")

	ownerToken, _ := ownerLogin("test@example.com", "hunter2")

	o, err := ownerGetByOwnerToken(ownerToken)
	if err != nil {
		t.Errorf("unexpected error on ownerGetByOwnerToken: %v", err)
		return
	}

	if o.OwnerHex != ownerHex {
		t.Errorf("expected ownerHex=%s got ownerHex=%s", ownerHex, o.OwnerHex)
		return
	}
}

func TestOwnerGetByOwnerTokenDNE(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if _, err := ownerGetByOwnerToken("does-not-exist"); err == nil {
		t.Errorf("expected error not found on ownerGetByOwnerToken before creating an account")
		return
	}
}
