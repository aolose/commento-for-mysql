package main

import (
	"fmt"
	"net/http"
	"os"
)

func ownerConfirmHex(confirmHex string) error {
	if confirmHex == "" {
		return errorMissingField
	}

	statement := `
		UPDATE owners
		SET confirmed_email=true
		WHERE owner_hex IN (
			SELECT owner_hex FROM owner_confirm_hexes
			WHERE confirm_hex=?
		);
	`
	res := db.Exec(statement, confirmHex)
	if res.Error != nil {
		logger.Errorf("cannot mark user's confirmedEmail as true: %v\n", res.Error)
		return errorInternal
	}

	count := res.RowsAffected

	if count == 0 {
		return errorNoSuchConfirmationToken
	}

	statement = `
		DELETE FROM owner_confirm_hexes
		WHERE confirm_hex=?;
	`
	err := db.Exec(statement, confirmHex).Error
	if err != nil {
		logger.Warningf("cannot remove confirmation token: %v\n", err)
		// Don't return an error because this is not critical.
	}

	return nil
}

func ownerConfirmHexHandler(w http.ResponseWriter, r *http.Request) {
	if confirmHex := r.FormValue("confirmHex"); confirmHex != "" {
		if err := ownerConfirmHex(confirmHex); err == nil {
			http.Redirect(w, r, fmt.Sprintf("%s/login?confirmed=true", os.Getenv("ORIGIN")), http.StatusTemporaryRedirect)
			return
		}
	}

	// TODO: include error message in the URL
	http.Redirect(w, r, fmt.Sprintf("%s/login?confirmed=false", os.Getenv("ORIGIN")), http.StatusTemporaryRedirect)
}
