package main

import (
	"database/sql"
	"log"
	"time"
)

type SqliteDatabase struct {
	*sql.DB
}

func sqliteInit(dbParams map[string]interface{}) (*SqliteDatabase, error) {

	dbFilename, ok := dbParams["file"].(string)
	if !ok {
		return nil, Error("err.conn.parse.sqlite.file.missing")
	}

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return nil, err
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
	if _, err = db.Exec(statement); err != nil {
		return nil, err
	} else {
		return &SqliteDatabase{db}, nil
	}
}

func (db *SqliteDatabase) CreateComment(comment *Comment) error {
	url := comment.URL
	name := comment.Name
	text := comment.Comment
	parent := comment.Parent

	statement := `
		SELECT depth, parent FROM comments WHERE rowid=?;
	`
	rows, err := db.Query(statement, parent)
	if err != nil {
		return err
	}
	defer rows.Close()

	depth := 0

	for rows.Next() {
		var pParent int
		if err := rows.Scan(&depth, &pParent); err == nil {
			if depth+1 > 5 {
				parent = pParent
			}
		}
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	statement = `
		INSERT INTO comments(url, name, text, time, depth, parent) VALUES(?, ?, ?, ?, ?, ?);
	`
	_, err = db.Exec(statement, url, name, text, time.Now(), depth+1, parent)
	return err
}

func (db *SqliteDatabase) GetComments(url string) ([]Comment, error) {
	statement := `
		SELECT rowid, url, comment, name, time, parent FROM comments WHERE url=?;
	`
	rows, err := db.Query(statement, url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var id int
		var url string
		var comment string
		var name string
		var parent int
		var timestamp time.Time
		if err = rows.Scan(&id, &url, &comment, &name, &timestamp, &parent); err != nil {
			return nil, err
		}
		comments = append(comments, Comment{ID: id, URL: url, Comment: comment, Name: name, Timestamp: timestamp, Parent: parent})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	return comments, nil
}
