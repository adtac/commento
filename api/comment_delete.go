package main

import (
	"net/http"
	"time"
)

func commentDelete(commentHex string, deleterHex string) error {
	if commentHex == "" || deleterHex == "" {
		return errorMissingField
	}

	statement := `
		UPDATE comments
		SET
			deleted = true,
			markdown = '[deleted]',
			html = '[deleted]',
			commenterHex = 'anonymous',
			deleterHex = $2,
			deletionDate = $3
		WHERE commentHex = $1;
	`
	_, err := db.Exec(statement, commentHex, deleterHex, time.Now().UTC())

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

	cm, err := commentGetByCommentHex(*x.CommentHex)
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

	if !isModerator && cm.CommenterHex != c.CommenterHex {
		bodyMarshal(w, response{"success": false, "message": errorNotModerator.Error()})
		return
	}

	if err = commentDelete(*x.CommentHex, c.CommenterHex); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
