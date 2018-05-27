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
		Session    *string `json:"session"`
		CommentHex *string `json:"commentHex"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetBySession(*x.Session)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	domain, err := commentDomainGet(*x.CommentHex)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	isModerator, err := isDomainModerator(c.Email, domain)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isModerator {
		writeBody(w, response{"success": false, "message": errorNotModerator.Error()})
		return
	}

	if err = commentApprove(*x.CommentHex); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true})
}
