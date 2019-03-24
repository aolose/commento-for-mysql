package main

import (
	"net/http"
)

func emailUpdate(e email) error {
	err := db.Model(&Emails{}).
		Where("email = ? AND unsubscribeSecretHex = ?", e.Email, e.UnsubscribeSecretHex).
		Update(Emails{SendModeratorNotifications: e.SendModeratorNotifications, SendReplyNotifications: e.SendReplyNotifications}).Error
	if err != nil {
		logger.Errorf("error updating email: %v", err)
		return errorInternal
	}

	return nil
}

func emailUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email *email `json:"email"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := emailUpdate(*x.Email); err != nil {
		bodyMarshal(w, response{"success": true, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
