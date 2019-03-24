package main

import (
	"net/http"
)

func commentList(commenterHex string, domain string, path string, includeUnapproved bool) ([]comment, map[string]commenter, error) {
	// path can be empty
	if commenterHex == "" || domain == "" {
		return nil, nil, errorMissingField
	}
	a := "domain = ? AND path = ? "
	if !includeUnapproved {
		if commenterHex == "anonymous" {
			a += `AND state = 'approved'`
		} else {
			a += "AND (state = 'approved' OR commenter_hex = '" + commenterHex + "')"
		}
	}

	rows, err := db.Model(&Comments{}).Where(a, domain, path).Rows()

	if err != nil {
		logger.Errorf("cannot get comments: %v", err)
		return nil, nil, errorInternal
	}
	defer rows.Close()

	commenters := make(map[string]commenter)
	commenters["anonymous"] = commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []comment{}
	for rows.Next() {
		c := comment{}
		if err = rows.Scan(&c.CommentHex, &c.Domain, &c.Path, &c.CommenterHex, &c.Markdown, &c.Html, &c.ParentHex, &c.Score, &c.State, &c.CreationDate); err != nil {
			logger.Errorf("%v", err)
			return nil, nil, errorInternal
		}
		if commenterHex != "anonymous" {
			row, err := db.Table("votes").Select("direction").
				Where("comment_hex = ? AND commenter_hex = ?", c.CommentHex, commenterHex).Rows()
			if err = row.Scan(&c.Direction); err != nil {
				c.Direction = 0
			}
		}

		if !includeUnapproved {
			c.State = ""
		}

		comments = append(comments, c)

		if _, ok := commenters[c.CommenterHex]; !ok {
			commenters[c.CommenterHex], err = commenterGetByHex(c.CommenterHex)
			if err != nil {
				logger.Errorf("cannot retrieve commenter: %v", err)
				return nil, nil, errorInternal
			}
		}
	}

	return comments, commenters, nil
}

func commentListHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"CommenterToken"`
		Domain         *string `json:"domain"`
		Path           *string `json:"path"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	path := *x.Path

	d, err := domainGet(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	p, err := pageGet(domain, path)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	commenterHex := "anonymous"
	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		if err == errorNoSuchToken {
			commenterHex = "anonymous"
		} else {
			bodyMarshal(w, response{"success": false, "message": err.Error()})
			return
		}
	} else {
		commenterHex = c.CommenterHex
	}

	isModerator := false
	modList := map[string]bool{}
	for _, mod := range d.Moderators {
		modList[mod.Email] = true
		if mod.Email == c.Email {
			isModerator = true
		}
	}

	domainViewRecord(domain, commenterHex)

	comments, commenters, err := commentList(commenterHex, domain, path, isModerator)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	_commenters := map[string]commenter{}
	for commenterHex, cr := range commenters {
		if _, ok := modList[cr.Email]; ok {
			cr.IsModerator = true
		}
		cr.Email = ""
		_commenters[commenterHex] = cr
	}

	bodyMarshal(w, response{
		"success":               true,
		"domain":                domain,
		"comments":              comments,
		"commenters":            _commenters,
		"requireModeration":     d.RequireModeration,
		"requireIdentification": d.RequireIdentification,
		"isFrozen":              d.State == "frozen",
		"isModerator":           isModerator,
		"attributes":            p,
		"configuredOauths":      configuredOauths,
	})
}
