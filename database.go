package main

import (
	"database/sql"
)

func loadDatabase(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	statement := `
		CREATE TABLE IF NOT EXISTS comments (
			url text not null,
			name text not null,
			comment text not null,
			time timestamp not null,
			parent int
		);
	`
	_, err = db.Exec(statement)
	return err

}

func cleanupOldComments() error {
	statement := `
		DELETE FROM comments
		WHERE time < date('now', '-30 minute');
	`
	_, err := db.Exec(statement)
	return err
}
