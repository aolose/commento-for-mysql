package main

import (
	"testing"
)

func TestCommenterSessionUpdateBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	commenterToken, _ := commenterTokenNew()

	if err := commenterSessionUpdate(commenterToken, "temp-commenter-hex"); err != nil {
		t.Errorf("unexpected error updating commenter session: %v", err)
		return
	}

	statement := `
		SELECT commenter_hex
		FROM commenter_sessions
		WHERE commenter_token = ?;
	`
	row := db.Raw(statement, commenterToken).Row()

	var commenterHex string
	if err := row.Scan(&commenterHex); err != nil {
		t.Errorf("error scanning commenterHex: %v", err)
		return
	}

	if commenterHex != "temp-commenter-hex" {
		t.Errorf("expected commenterHex=temp-commenter-hex got commenterHex=%s", commenterHex)
		return
	}
}

func TestCommenterSessionUpdateEmpty(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	if err := commenterSessionUpdate("", "temp-commenter-hex"); err == nil {
		t.Errorf("expected error not found when updating with empty commenterToken")
		return
	}
}
