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
	statement := `
		INSERT INTO
		pages  (domain, path, isLocked, stickyCommentHex)
		VALUES ($1,     $2,   $3,       $4              )
		ON CONFLICT (domain, path) DO
			UPDATE SET isLocked = $3, stickyCommentHex = $4;
	`
	_, err := db.Exec(statement, p.Domain, p.Path, p.IsLocked, p.StickyCommentHex)
	if err != nil {
		logger.Errorf("error setting page attributes: %v", err)
		return errorInternal
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
