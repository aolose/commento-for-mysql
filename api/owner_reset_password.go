package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ownerResetPassword(resetHex string, password string) error {
	if resetHex == "" || password == "" {
		return errorMissingField
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("cannot generate hash from password: %v\n", err)
		return errorInternal
	}

	statement := `
		UPDATE owners SET password_hash = ?
		WHERE owner_hex = (
			SELECT owner_hex
			FROM owner_resetHexes
			WHERE reset_hex=?
		);
	`
	res := db.Exec(statement, string(passwordHash), resetHex)
	if res.Error != nil {
		logger.Errorf("cannot change user's password: %v\n", res.Error)
		return errorInternal
	}

	count := res.RowsAffected

	if count == 0 {
		return errorNoSuchResetToken
	}

	statement = `
		DELETE FROM owner_resetHexes
    WHERE reset_hex=?;
	`
	err = db.Exec(statement, resetHex).Error
	if err != nil {
		logger.Warningf("cannot remove reset token: %v\n", err)
	}

	return nil
}

func ownerResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ResetHex *string `json:"resetHex"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := ownerResetPassword(*x.ResetHex, *x.Password); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
