package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func LoadDatabase(dbString string) error {
	var err error
	db, err := sql.Open("postgres", dbString)
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

func CleanupOldComments() error {
	statement := `
		DELETE FROM comments
		WHERE time < date('now', '-30 minute');
	`
	_, err := db.Exec(statement)
	return err
}
