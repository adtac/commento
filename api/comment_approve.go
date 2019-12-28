package main

import (
	"net/http"
)

func commentApprove(commentHex string) error {
	if commentHex == "" {
		return errorMissingField
	}

	statement := `
		UPDATE comments
		SET state = 'approved'
		WHERE commentHex = $1;
	`

	_, err := db.Exec(statement, commentHex)
	if err != nil {
		logger.Errorf("cannot approve comment: %v", err)
		return errorInternal
	}

	return nil
}

func commentApproveHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = commentApprove(*x.CommentHex); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
