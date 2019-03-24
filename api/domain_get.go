package main

import ()

func domainGet(dmn string) (domain, error) {
	if dmn == "" {
		return domain{}, errorMissingField
	}

	statement := `
    SELECT domain, owner_hex, name, creation_date, state, imported_comments, auto_spam_filter, require_moderation, require_identification, moderate_all_anonymous, email_notification_policy
		FROM domains
		WHERE domain = ?;
	`
	row := db.Raw(statement, dmn).Row()

	var err error
	d := domain{}
	if err = row.Scan(&d.Domain, &d.OwnerHex, &d.Name, &d.CreationDate, &d.State, &d.ImportedComments, &d.AutoSpamFilter, &d.RequireModeration, &d.RequireIdentification, &d.ModerateAllAnonymous, &d.EmailNotificationPolicy); err != nil {
		return d, errorNoSuchDomain
	}

	d.Moderators, err = domainModeratorList(d.Domain)
	if err != nil {
		return domain{}, err
	}

	return d, nil
}
