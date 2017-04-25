package main

import (
	"database/sql"
)

func createTables() error {
	statement := `
		CREATE TABLE comments (
			url text not null,
			name text not null,
			comment text not null,
			time timestamp not null,
			parent int
		);
	`
	_, err := db.Exec(statement)
	return err
}

func loadDatabase(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	statement := `
		SELECT name FROM sqlite_master WHERE type='table' AND name='comments';
	`
	rows, err := db.Query(statement)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		if err = createTables(); err != nil {
			return err
		}
	}

	return nil
}
