package main

import (
	"net/http"
)

func pageUpdate(p page) error {
	if p.Domain == "" {
		return errorMissingField
	}

	// fields to not update:
	//   commentCount
	err := db.Where(Pages{Domain: p.Domain, Path: p.Path}).
		Assign(Pages{IsLocked: p.IsLocked, StickyCommentHex: p.StickyCommentHex}).
		FirstOrCreate(&Pages{}).Error

	if err != nil {
		logger.Errorf("error setting page attributes: %v", err)
		//return errorInternal
	}

	return nil
}

func pageUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		Domain         *string `json:"domain"`
		Path           *string `json:"path"`
		Attributes     *page   `json:"attributes"`
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

	domain := domainStrip(*x.Domain)

	isModerator, err := isDomainModerator(domain, c.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isModerator {
		bodyMarshal(w, response{"success": false, "message": errorNotModerator.Error()})
		return
	}

	(*x.Attributes).Domain = *x.Domain
	(*x.Attributes).Path = *x.Path

	if err = pageUpdate(*x.Attributes); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
