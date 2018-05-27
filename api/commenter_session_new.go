package main

import (
	"net/http"
	"time"
)

func commenterSessionNew() (string, error) {
	session, err := randomHex(32)
	if err != nil {
		logger.Errorf("cannot create session hex: %v", err)
		return "", errorInternal
	}

	statement := `
		INSERT INTO
		commenterSessions (session, creationDate)
		VALUES            ($1,      $2          );
	`
	_, err = db.Exec(statement, session, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert new session: %v", err)
		return "", errorInternal
	}

	return session, nil
}

func commenterSessionNewHandler(w http.ResponseWriter, r *http.Request) {
	session, err := commenterSessionNew()
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true, "session": session})
}
