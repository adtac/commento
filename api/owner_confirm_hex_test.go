package main

import (
	"testing"
	"time"
)

func TestOwnerConfirmHexBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	ownerHex, _ := ownerNew("test@example.com", "Test", "hunter2")

	statement := `
		UPDATE owners
		SET confirmedEmail=false;
	`
	_, err := db.Exec(statement)
	if err != nil {
		t.Errorf("unexpected error when setting confirmedEmail=false: %v", err)
		return
	}

	confirmHex, _ := randomHex(32)

	statement = `
		INSERT INTO
		ownerConfirmHexes (confirmHex, ownerHex, sendDate)
		VALUES            ($1,         $2,       $3      );
	`
	_, err = db.Exec(statement, confirmHex, ownerHex, time.Now().UTC())
	if err != nil {
		t.Errorf("unexpected error creating inserting confirmHex: %v\n", err)
		return
	}

	if err = ownerConfirmHex(confirmHex); err != nil {
		t.Errorf("unexpected error confirming hex: %v", err)
		return
	}

	statement = `
		SELECT confirmedEmail
		FROM owners
		WHERE ownerHex=$1;
	`
	row := db.QueryRow(statement, ownerHex)

	var confirmedHex bool
	if err = row.Scan(&confirmedHex); err != nil {
		t.Errorf("unexpected error scanning confirmedEmail: %v", err)
		return
	}

	if !confirmedHex {
		t.Errorf("confirmedHex expected to be true after confirmation; found to be false")
		return
	}
}
