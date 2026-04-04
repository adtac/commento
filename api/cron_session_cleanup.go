package main

import (
	"time"
)

func sessionCleanupBegin() error {
	go func() {
		for {
			sessionCleanup()
			time.Sleep(1 * time.Hour)
		}
	}()

	return nil
}

func sessionCleanup() {
	statement := `
		DELETE FROM ownerSessions
		WHERE expiresAt < NOW();
	`
	result, err := db.Exec(statement)
	if err != nil {
		logger.Warningf("cannot clean up expired sessions: %v", err)
		return
	}

	count, _ := result.RowsAffected()
	if count > 0 {
		logger.Infof("cleaned up %d expired owner sessions", count)
	}
}
