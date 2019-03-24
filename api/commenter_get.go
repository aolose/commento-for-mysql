package main

func commenterGetByHex(commenterHex string) (commenter, error) {
	if commenterHex == "" {
		return commenter{}, errorMissingField
	}

	statement := `
    SELECT commenter_hex, email, name, link, photo, provider, join_date
    FROM commenters
    WHERE commenter_hex = ?;
  `
	row := db.Raw(statement, commenterHex).Row()

	c := commenter{}
	if err := row.Scan(&c.CommenterHex, &c.Email, &c.Name, &c.Link, &c.Photo, &c.Provider, &c.JoinDate); err != nil {
		// TODO: is this the only error?
		return commenter{}, errorNoSuchCommenter
	}

	return c, nil
}

func commenterGetByEmail(provider string, email string) (commenter, error) {
	if provider == "" || email == "" {
		return commenter{}, errorMissingField
	}

	statement := `
    SELECT commenter_hex, email, name, link, photo, provider, join_date
    FROM commenters
    WHERE email = ? AND provider = ?;
  `
	row := db.Raw(statement, email, provider).Row()

	c := commenter{}
	if err := row.Scan(&c.CommenterHex, &c.Email, &c.Name, &c.Link, &c.Photo, &c.Provider, &c.JoinDate); err != nil {
		// TODO: is this the only error?
		return commenter{}, errorNoSuchCommenter
	}

	return c, nil
}

func commenterGetByCommenterToken(commenterToken string) (commenter, error) {
	if commenterToken == "" {
		return commenter{}, errorMissingField
	}

	statement := `
    SELECT commenter_hex
    FROM commenter_sessions
    WHERE commenter_token = ?;
	`
	row := db.Raw(statement, commenterToken).Row()

	var commenterHex string
	if err := row.Scan(&commenterHex); err != nil {
		// TODO: is the only error?
		return commenter{}, errorNoSuchToken
	}

	if commenterHex == "none" {
		return commenter{}, errorNoSuchToken
	}

	// TODO: use a join instead of two queries?
	return commenterGetByHex(commenterHex)
}
