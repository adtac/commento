package main

import ()

func commentOwnershipVerify(commenterHex string, commentHex string) (bool, error) {
	if commenterHex == "" || commentHex == "" {
		return false, errorMissingField
	}

	statement := `
		SELECT EXISTS (
			SELECT 1
			FROM comments
			WHERE commenterHex=$1 AND commentHex=$2
		);
	`
	row := db.QueryRow(statement, commenterHex, commentHex)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		logger.Errorf("cannot query if comment owner: %v", err)
		return false, errorInternal
	}

	return exists, nil
}
