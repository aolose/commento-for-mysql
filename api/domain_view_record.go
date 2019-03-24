package main

import (
	"time"
)

func domainViewRecord(domain string, commenterHex string) {
	statement := `
		INSERT INTO
		views  (domain, commenter_hex, view_date)
		VALUES (?,     ?,           ?      );
	`
	err := db.Exec(statement, domain, commenterHex, time.Now().UTC()).Error
	if err != nil {
		logger.Warningf("cannot insert views: %v", err)
	}
}
