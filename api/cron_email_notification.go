package main

import (
	"time"
)

func emailNotificationBegin() error {
	go func() {
		for {
			statement := `
				SELECT email, send_moderator_notifications, send_reply_notifications
				FROM emails
				WHERE pending_emails > 0 AND last_email_notification_date < ?;
			`
			rows, err := db.Raw(statement, time.Now().UTC().Add(time.Duration(-10)*time.Minute)).Rows()
			if err != nil {
				logger.Errorf("cannot query domains: %v", err)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var email string
				var sendModeratorNotifications bool
				var sendReplyNotifications bool
				if err = rows.Scan(&email, &sendModeratorNotifications, &sendReplyNotifications); err != nil {
					logger.Errorf("cannot scan email in cron job to send notifications: %v", err)
					continue
				}

				if _, ok := emailQueue[email]; !ok {
					if err = emailNotificationPendingReset(email); err != nil {
						logger.Errorf("error resetting pendingEmails: %v", err)
						continue
					}
				}

				cont := true
				kindListMap := map[string][]emailNotification{}
				for cont {
					select {
					case e := <-emailQueue[email]:
						if _, ok := kindListMap[e.Kind]; !ok {
							kindListMap[e.Kind] = []emailNotification{}
						}

						if (e.Kind == "reply" && sendReplyNotifications) || sendModeratorNotifications {
							kindListMap[e.Kind] = append(kindListMap[e.Kind], e)
						}
					default:
						cont = false
						break
					}
				}

				for kind, list := range kindListMap {
					go emailNotificationSend(email, kind, list)
				}
			}

			time.Sleep(10 * time.Minute)
		}
	}()

	return nil
}
