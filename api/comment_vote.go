package main

import (
	"net/http"
	"time"
)

func commentVote(commenterHex string, commentHex string, direction int) error {
	if commentHex == "" || commenterHex == "" {
		return errorMissingField
	}

	statement := `
		SELECT commenterHex
		FROM comments
		WHERE commentHex = $1;
	`
	row := db.QueryRow(statement, commentHex)

	var authorHex string
	if err := row.Scan(&authorHex); err != nil {
		logger.Errorf("error selecting authorHex for vote")
		return errorInternal
	}

	if authorHex == commenterHex {
		return errorSelfVote
	}

	statement = `
		INSERT INTO
		votes  (commentHex, commenterHex, direction, voteDate)
		VALUES ($1,         $2,           $3,        $4      )
		ON CONFLICT (commentHex, commenterHex) DO
		UPDATE SET direction = $3;
	`
	_, err := db.Exec(statement, commentHex, commenterHex, direction, time.Now().UTC())
	if err != nil {
		logger.Errorf("error inserting/updating votes: %v", err)
		return errorInternal
	}

	return nil
}

func commentVoteHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		CommentHex     *string `json:"commentHex"`
		Direction      *int    `json:"direction"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if *x.CommenterToken == "anonymous" {
		bodyMarshal(w, response{"success": false, "message": errorUnauthorisedVote.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	direction := 0
	if *x.Direction > 0 {
		direction = 1
	} else if *x.Direction < 0 {
		direction = -1
	}

	if err := commentVote(c.CommenterHex, *x.CommentHex, direction); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
