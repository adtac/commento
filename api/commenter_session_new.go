package main

import (
	"net/http"
	"time"
)

func commenterTokenNew() (string, error) {
	commenterToken, err := randomHex(32)
	if err != nil {
		logger.Errorf("cannot create commenterToken: %v", err)
		return "", errorInternal
	}

	statement := `
		INSERT INTO
		commenterSessions (commenterToken, creationDate)
		VALUES            ($1,             $2          );
	`
	_, err = db.Exec(statement, commenterToken, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert new commenterToken: %v", err)
		return "", errorInternal
	}

	return commenterToken, nil
}

func commenterTokenNewHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken, err := commenterTokenNew()
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commenterToken": commenterToken})
}
