package main

import ()

func domainOwnershipVerify(ownerHex string, domain string) (bool, error) {
	if ownerHex == "" || domain == "" {
		return false, errorMissingField
	}

	statement := `
		SELECT EXISTS (
			SELECT 1
			FROM domains
			WHERE owner_hex=? AND domain=?
		);
	`
	row := db.Raw(statement, ownerHex, domain).Row()

	var exists bool
	if err := row.Scan(&exists); err != nil {
		logger.Errorf("cannot query if domain owner: %v", err)
		return false, errorInternal
	}

	return exists, nil
}
