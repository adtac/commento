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
			url varchar(2083) not null,
			name varchar(200) not null,
			comment varchar(3000) not null,
			depth int not null,
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
