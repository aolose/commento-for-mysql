package main

import (
	"net/http"
)

func commentDelete(commentHex string) error {
	if commentHex == "" {
		return errorMissingField
	}

	err := db.Delete(&Comments{}, "comment_hex = ?", commentHex).Error

	if err != nil {
		// TODO: make sure this is the error is actually non-existant commentHex
		return errorNoSuchComment
	}

	return nil
}

func commentDeleteHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		CommentHex     *string `json:"commentHex"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain, _, err := commentDomainPathGet(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	isModerator, err := isDomainModerator(domain, c.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isModerator {
		bodyMarshal(w, response{"success": false, "message": errorNotModerator.Error()})
		return
	}

	if err = commentDelete(*x.CommentHex); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
