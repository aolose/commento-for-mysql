package main

import ()

func ownerGetByEmail(email string) (owner, error) {
	if email == "" {
		return owner{}, errorMissingField
	}

	statement := `
    SELECT owner_hex, email, name, confirmed_email, join_date
    FROM owners
    WHERE email=?;
  `
	row := db.Raw(statement, email).Row()

	var o owner
	if err := row.Scan(&o.OwnerHex, &o.Email, &o.Name, &o.ConfirmedEmail, &o.JoinDate); err != nil {
		// TODO: Make sure this is actually no such email.
		return owner{}, errorNoSuchEmail
	}

	return o, nil
}

func ownerGetByOwnerToken(ownerToken string) (owner, error) {
	if ownerToken == "" {
		return owner{}, errorMissingField
	}

	statement := `
    SELECT owner_hex, email, name, confirmed_email, join_date
		FROM owners
		WHERE owner_hex IN (
			SELECT owner_hex FROM owner_sessions
			WHERE owner_token = ?
		);
	`
	row := db.Raw(statement, ownerToken).Row()

	var o owner
	if err := row.Scan(&o.OwnerHex, &o.Email, &o.Name, &o.ConfirmedEmail, &o.JoinDate); err != nil {
		logger.Errorf("cannot scan owner: %v\n", err)
		return owner{}, errorInternal
	}

	return o, nil
}

func ownerGetByOwnerHex(ownerHex string) (owner, error) {
	if ownerHex == "" {
		return owner{}, errorMissingField
	}

	statement := `
    SELECT owner_hex, email, name, confirmed_email, join_date
		FROM owners
		WHERE owner_hex = ?;
	`
	row := db.Raw(statement, ownerHex).Row()

	var o owner
	if err := row.Scan(&o.OwnerHex, &o.Email, &o.Name, &o.ConfirmedEmail, &o.JoinDate); err != nil {
		logger.Errorf("cannot scan owner: %v\n", err)
		return owner{}, errorInternal
	}

	return o, nil
}
