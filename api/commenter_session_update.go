package main

import ()

func commenterSessionUpdate(session string, commenterHex string) error {
	if session == "" || commenterHex == "" {
		return errorMissingField
	}

	statement := `
    UPDATE commenterSessions
    SET commenterHex=$2
    WHERE session=$1;
  `
	_, err := db.Exec(statement, session, commenterHex)
	if err != nil {
		logger.Errorf("error updating commenterHex in commenterSessions: %v", err)
		return errorInternal
	}

	return nil
}
