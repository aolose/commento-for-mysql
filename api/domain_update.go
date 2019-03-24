package main

import (
	"net/http"
)

func domainUpdate(d domain) error {
	statement := `
		UPDATE domains
    SET name=?, state=?, auto_spam_filter=?, require_moderation=?, require_identification=?, moderate_all_anonymous=?, email_notification_policy=?
		WHERE domain=?;
	`

	err := db.Exec(statement, d.Name, d.State, d.AutoSpamFilter, d.RequireModeration, d.RequireIdentification, d.ModerateAllAnonymous, d.EmailNotificationPolicy, d.Domain).Error
	if err != nil {
		logger.Errorf("cannot update non-moderators: %v", err)
		return errorInternal
	}

	return nil
}

func domainUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		D          *domain `json:"domain"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip((*x.D).Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err = domainUpdate(*x.D); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
