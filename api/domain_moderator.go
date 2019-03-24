package main

import (
	"time"
)

type moderator struct {
	Email   string    `json:"email"`
	Domain  string    `json:"domain"`
	AddDate time.Time `json:"addDate"`
}

func domainModeratorList(domain string) ([]moderator, error) {
	statement := `
		SELECT email, add_date
		FROM moderators
		WHERE domain=?;
	`
	rows, err := db.Raw(statement, domain).Rows()
	if err != nil {
		logger.Errorf("cannot get moderators: %v", err)
		return nil, errorInternal
	}
	defer rows.Close()

	moderators := []moderator{}
	for rows.Next() {
		m := moderator{}
		if err = rows.Scan(&m.Email, &m.AddDate); err != nil {
			logger.Errorf("cannot Scan moderator: %v", err)
			return nil, errorInternal
		}

		moderators = append(moderators, m)
	}

	return moderators, nil
}

func isDomainModerator(domain string, email string) (bool, error) {
	statement := `
		SELECT EXISTS (
			SELECT 1
			FROM moderators
			WHERE domain=? AND email=?
		);
	`
	row := db.Raw(statement, domain, email).Row()

	var exists bool
	if err := row.Scan(&exists); err != nil {
		logger.Errorf("cannot query if moderator: %v", err)
		return false, errorInternal
	}

	return exists, nil
}
