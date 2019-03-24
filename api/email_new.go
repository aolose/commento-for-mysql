package main

import (
	"time"
)

func emailNew(email string) error {
	unsubscribeSecretHex, err := randomHex(32)
	if err != nil {
		logger.Errorf("%v", err)
		return errorInternal
	}
	var m Emails
	db.First(&m, "email = ?", email).Scan(&m)
	if m.Email == "" {
		err = db.Create(&Emails{Email: email, UnsubscribeSecretHex: unsubscribeSecretHex, LastEmailNotificationDate: time.Now().UTC()}).
			Error
		if err != nil {
			logger.Errorf("cannot insert email into emails: %v", err)
			return errorInternal
		}
	}
	return nil
}
