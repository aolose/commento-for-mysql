package main

import (
	"html/template"
)

func emailNotificationSend(em string, kind string, notifications []emailNotification) {
	if len(notifications) == 0 {
		return
	}

	e, err := emailGet(em)
	if err != nil {
		logger.Errorf("cannot get email: %v", err)
		return
	}

	messages := []emailNotificationText{}
	for _, notification := range notifications {
		statement := `
			SELECT html
			FROM comments
			WHERE comment_hex = ?;
		`
		row := db.Raw(statement, notification.CommentHex).Row()

		var html string
		if err = row.Scan(&html); err != nil {
			// the comment was deleted?
			// TODO: is this the only error?
			return
		}

		messages = append(messages, emailNotificationText{
			emailNotification: notification,
			Html:              template.HTML(html),
		})
	}

	statement := `
		SELECT name
		FROM commenters
		WHERE email = ?;
	`
	row := db.Raw(statement, em).Row()

	var name string
	if err := row.Scan(&name); err != nil {
		// The moderator has probably not created a commenter account. Let's just
		// use their email as name.
		name = nameFromEmail(em)
	}

	if err := emailNotificationPendingReset(em); err != nil {
		logger.Errorf("cannot reset after email notification: %v", err)
		return
	}

	if err := smtpEmailNotification(em, name, e.UnsubscribeSecretHex, messages, kind); err != nil {
		logger.Errorf("cannot send email notification: %v", err)
		return
	}
}
