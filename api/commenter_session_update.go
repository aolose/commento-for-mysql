package main

import ()

func commenterSessionUpdate(commenterToken string, commenterHex string) error {
	if commenterToken == "" || commenterHex == "" {
		return errorMissingField
	}

	statement := `
    UPDATE commenter_sessions
    SET commenter_hex = ?
    WHERE commenter_token = ?;
  `
	err := db.Exec(statement, commenterHex, commenterToken).Error
	if err != nil {
		logger.Errorf("error updating commenterHex: %v", err)
		return errorInternal
	}

	return nil
}
