package main

import (
	"database/sql"
)

var db *sql.DB

// LoadDatabase connects to the commento database
func LoadDatabase(dbFile string) error {
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

// CleanupOldComments deletes old comments, useful in demo environments
func CleanupOldComments() error {
	statement := `
		DELETE FROM comments
		WHERE time < date('now', '-30 minute');
	`
	_, err := db.Exec(statement)
	return err
}
