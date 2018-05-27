package main

import (
	"testing"
	"time"
)

func TestOwnerResetPasswordBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	ownerHex, _ := ownerNew("test@example.com", "Test", "hunter2")

	resetHex, _ := randomHex(32)

	statement := `
		INSERT INTO
		ownerResetHexes (resetHex, ownerHex, sendDate)
		VALUES          ($1,       $2,    $3         );
	`
	_, err := db.Exec(statement, resetHex, ownerHex, time.Now().UTC())
	if err != nil {
		t.Errorf("unexpected error inserting resetHex: %v", err)
		return
	}

	if err = ownerResetPassword(resetHex, "hunter3"); err != nil {
		t.Errorf("unexpected error resetting password: %v", err)
		return
	}

	if _, err := ownerLogin("test@example.com", "hunter2"); err == nil {
		t.Errorf("expected error not found when given old password")
		return
	}

	if _, err := ownerLogin("test@example.com", "hunter3"); err != nil {
		t.Errorf("unexpected error when logging in: %v", err)
		return
	}
}
