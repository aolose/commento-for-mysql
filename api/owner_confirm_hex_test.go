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
    SET confirmed_email=false;
  `
	err := db.Exec(statement).Error
	if err != nil {
		t.Errorf("unexpected error when setting confirmedEmail=false: %v", err)
		return
	}

	confirmHex, _ := randomHex(32)

	statement = `
    INSERT INTO
    owner_confirm_hexes (confirm_hex, owner_hex, send_date)
    VALUES            (?,         ?,       ?      );
  `
	err = db.Exec(statement, confirmHex, ownerHex, time.Now().UTC()).Error
	if err != nil {
		t.Errorf("unexpected error creating inserting confirmHex: %v\n", err)
		return
	}

	if err = ownerConfirmHex(confirmHex); err != nil {
		t.Errorf("unexpected error confirming hex: %v", err)
		return
	}

	statement = `
    SELECT confirmed_email
    FROM owners
    WHERE owner_hex=?;
  `
	row := db.Raw(statement, ownerHex).Row()

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
