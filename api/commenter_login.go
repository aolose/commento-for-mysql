package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func commenterLogin(email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", errorMissingField
	}

	statement := `
		SELECT commenter_hex, password_hash
		FROM commenters
		WHERE email = ? AND provider = 'commento';
	`
	row := db.Table("commenters").
		Select("commenter_hex, password_hash").
		Where("email = ? AND provider = 'commento'", email).Row()

	var commenterHex string
	var passwordHash string
	if err := row.Scan(&commenterHex, &passwordHash); err != nil {
		return "", errorInvalidEmailPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		// TODO: is this the only possible error?
		return "", errorInvalidEmailPassword
	}

	commenterToken, err := randomHex(32)
	if err != nil {
		logger.Errorf("cannot create commenterToken: %v", err)
		return "", errorInternal
	}

	statement = `
		INSERT INTO
		commenter_sessions (commenter_token, commenter_hex, creation_date)
		VALUES            (?,             ?,           ?          );
	`
	err = db.Exec(statement, commenterToken, commenterHex, time.Now().UTC()).Error
	if err != nil {
		logger.Errorf("cannot insert commenterToken token: %v\n", err)
		return "", errorInternal
	}

	return commenterToken, nil
}

func commenterLoginHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	commenterToken, err := commenterLogin(*x.Email, *x.Password)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// TODO: modify commenterLogin to directly return c?
	c, err := commenterGetByCommenterToken(commenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commenterToken": commenterToken, "commenter": c})
}
