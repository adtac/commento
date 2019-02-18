package main

import (
	"time"
)

type emailNotification struct {
	Email         string
	CommenterName string
	Domain        string
	Path          string
	Title         string
	CommentHex    string
	Kind          string
}

var emailQueue map[string](chan emailNotification) = map[string](chan emailNotification){}

func emailNotificationPendingResetAll() error {
	statement := `
		UPDATE emails
		SET pendingEmails = 0;
	`
	_, err := db.Exec(statement)
	if err != nil {
		logger.Errorf("cannot reset pendingEmails: %v", err)
		return err
	}

	return nil
}

func emailNotificationPendingIncrement(email string) error {
	statement := `
		UPDATE emails
		SET pendingEmails = pendingEmails + 1
		WHERE email = $1;
	`
	_, err := db.Exec(statement, email)
	if err != nil {
		logger.Errorf("cannot increment pendingEmails: %v", err)
		return err
	}

	return nil
}

func emailNotificationPendingReset(email string) error {
	statement := `
		UPDATE emails
		SET pendingEmails = 0, lastEmailNotificationDate = $2
		WHERE email = $1;
	`
	_, err := db.Exec(statement, email, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot decrement pendingEmails: %v", err)
		return err
	}

	return nil
}

func emailNotificationEnqueue(e emailNotification) error {
	if err := emailNotificationPendingIncrement(e.Email); err != nil {
		logger.Errorf("cannot increment pendingEmails when enqueueing: %v", err)
		return err
	}

	if _, ok := emailQueue[e.Email]; !ok {
		// don't enqueue more than 10 emails as we won't send more than 10 comments
		// in one email anyway
		emailQueue[e.Email] = make(chan emailNotification, 10)
	}

	select {
	case emailQueue[e.Email] <- e:
	default:
	}

	return nil
}
