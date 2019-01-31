package main

import (
	"time"
)

func viewsCleanupBegin() error {
	go func() {
		for {
			statement := `
				DELETE FROM views
				WHERE viewDate < $1;
			`
			_, err := db.Exec(statement, time.Now().UTC().AddDate(0, 0, -45))
			if err != nil {
				logger.Errorf("error cleaning up views: %v", err)
				return
			}

			time.Sleep(24 * time.Hour)
		}
	}()

	return nil
}
