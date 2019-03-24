package main

import (
	"net/http"
	"time"
)

func commenterTokenNew() (string, error) {
	commenterToken, err := randomHex(32)
	if err != nil {
		logger.Errorf("cannot create commenterToken: %v", err)
		return "", errorInternal
	}

	statement := `
		INSERT INTO
		commenter_sessions (commenter_token, creation_date)
		VALUES            (?,             ?          );
	`
	err = db.Exec(statement, commenterToken, time.Now().UTC()).Error
	if err != nil {
		logger.Errorf("cannot insert new commenterToken: %v", err)
		return "", errorInternal
	}

	return commenterToken, nil
}

func commenterTokenNewHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken, err := commenterTokenNew()
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commenterToken": commenterToken})
}
