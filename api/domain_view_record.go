package main

import (
	"time"
)

func domainViewRecord(domain string, commenterHex string) {
	statement := `
		INSERT INTO
		views  (domain, commenterHex, viewDate)
		VALUES ($1,     $2,           $3      );
	`
	_, err := db.Exec(statement, domain, commenterHex, time.Now().UTC())
	if err != nil {
		logger.Warningf("cannot insert views: %v", err)
	}
}
