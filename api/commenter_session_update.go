package main

import ()

func commenterSessionUpdate(commenterToken string, commenterHex string) error {
	if commenterToken == "" || commenterHex == "" {
		return errorMissingField
	}

	statement := `
		UPDATE commenterSessions
		SET commenterHex = $2
		WHERE commenterToken = $1;
	`
	_, err := db.Exec(statement, commenterToken, commenterHex)
	if err != nil {
		logger.Errorf("error updating commenterHex: %v", err)
		return errorInternal
	}

	return nil
}
