package main

import ()

func commentOwnershipVerify(commenterHex string, commentHex string) (bool, error) {
	if commenterHex == "" || commentHex == "" {
		return false, errorMissingField
	}

	statement := `
		SELECT EXISTS (
			SELECT 1
			FROM comments
			WHERE commenter_hex=? AND comment_hex=?
		);
	`
	row := db.Raw(statement, commenterHex, commentHex).Row()

	var exists bool
	if err := row.Scan(&exists); err != nil {
		logger.Errorf("cannot query if comment owner: %v", err)
		return false, errorInternal
	}

	return exists, nil
}
