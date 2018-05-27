package main

import ()

func commenterSessionGet(session string) (commenterSession, error) {
	if session == "" {
		return commenterSession{}, errorMissingField
	}

	statement := `
    SELECT commenterHex, creationDate
    FROM commenterSessions
    WHERE session=$1;
  `
	row := db.QueryRow(statement, session)

	cs := commenterSession{}
	if err := row.Scan(&cs.CommenterHex, &cs.CreationDate); err != nil {
		return commenterSession{}, errorNoSuchSession
	}

	cs.Session = session

	return cs, nil
}
