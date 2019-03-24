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
		owner_reset_hexes (reset_hex, owner_hex, send_date)
		VALUES          (?,       ?,    ?         );
	`
	err := db.Exec(statement, resetHex, ownerHex, time.Now().UTC()).Error
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
