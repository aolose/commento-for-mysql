package main

import (
	"time"
)

type emailNotification struct {
	Email         string
	CommenterName string
	Domain        string
	Path          string
	Title         string
	CommentHex    string
	Kind          string
}

var emailQueue map[string](chan emailNotification) = map[string](chan emailNotification){}

func emailNotificationPendingResetAll() error {
	err := db.Model(&Emails{}).Updates(Emails{PendingEmails: 0}).Error
	if err != nil {
		logger.Errorf("cannot reset pendingEmails: %v", err)
		return err
	}
	return nil
}

func emailNotificationPendingIncrement(email string) error {
	statement := `
		UPDATE emails
		SET pending_emails = pending_emails + 1
		WHERE email = ?;
	`
	err := db.Exec(statement, email).Error
	if err != nil {
		logger.Errorf("cannot increment pendingEmails: %v", err)
		return err
	}

	return nil
}

func emailNotificationPendingReset(email string) error {
	statement := `
		UPDATE emails
		SET pending_emails = 0, last_email_notification_date = ?
		WHERE email = ?;
	`
	err := db.Exec(statement, time.Now().UTC(), email).Error
	if err != nil {
		logger.Errorf("cannot decrement pendingEmails: %v", err)
		return err
	}

	return nil
}

func emailNotificationEnqueue(e emailNotification) error {
	if err := emailNotificationPendingIncrement(e.Email); err != nil {
		logger.Errorf("cannot increment pendingEmails when enqueueing: %v", err)
		return err
	}

	if _, ok := emailQueue[e.Email]; !ok {
		// don't enqueue more than 10 emails as we won't send more than 10 comments
		// in one email anyway
		emailQueue[e.Email] = make(chan emailNotification, 10)
	}

	select {
	case emailQueue[e.Email] <- e:
	default:
	}

	return nil
}
