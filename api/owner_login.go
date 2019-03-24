package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func ownerLogin(email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", errorMissingField
	}

	statement := `
		SELECT owner_hex, confirmed_email, password_hash
		FROM owners
		WHERE email=?;
	`
	row := db.Raw(statement, email).Row()

	var ownerHex string
	var confirmedEmail bool
	var passwordHash string
	if err := row.Scan(&ownerHex, &confirmedEmail, &passwordHash); err != nil {
		return "", errorInvalidEmailPassword
	}

	if !confirmedEmail {
		return "", errorUnconfirmedEmail
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		// TODO: is this the only possible error?
		return "", errorInvalidEmailPassword
	}

	ownerToken, err := randomHex(32)
	if err != nil {
		logger.Errorf("cannot create ownerToken: %v", err)
		return "", errorInternal
	}

	statement = `
		INSERT INTO
		owner_sessions (owner_token, owner_hex, login_date)
		VALUES        (?,         ?,       ?       );
	`
	err = db.Exec(statement, ownerToken, ownerHex, time.Now().UTC()).Error
	if err != nil {
		logger.Errorf("cannot insert ownerSession: %v\n", err)
		return "", errorInternal
	}

	return ownerToken, nil
}

func ownerLoginHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	ownerToken, err := ownerLogin(*x.Email, *x.Password)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "ownerToken": ownerToken})
}
