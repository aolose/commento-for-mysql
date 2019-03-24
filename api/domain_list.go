package main

import (
	"net/http"
)

func domainList(ownerHex string) ([]domain, error) {
	if ownerHex == "" {
		return []domain{}, errorMissingField
	}

	statement := `
    SELECT domain, owner_hex, name, creation_date, state, imported_comments, auto_spam_filter, require_moderation, require_identification, moderate_all_anonymous, email_notification_policy
		FROM domains
		WHERE owner_hex=?;
	`
	rows, err := db.Raw(statement, ownerHex).Rows()
	if err != nil {
		logger.Errorf("cannot query domains: %v", err)
		return nil, errorInternal
	}
	defer rows.Close()

	domains := []domain{}
	for rows.Next() {
		d := domain{}
		if err = rows.Scan(&d.Domain, &d.OwnerHex, &d.Name, &d.CreationDate, &d.State, &d.ImportedComments, &d.AutoSpamFilter, &d.RequireModeration, &d.RequireIdentification, &d.ModerateAllAnonymous, &d.EmailNotificationPolicy); err != nil {
			logger.Errorf("cannot Scan domain: %v", err)
			return nil, errorInternal
		}

		d.Moderators, err = domainModeratorList(d.Domain)
		if err != nil {
			return []domain{}, err
		}

		domains = append(domains, d)
	}

	return domains, rows.Err()
}

func domainListHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
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

	domains, err := domainList(o.OwnerHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "domains": domains})
}
