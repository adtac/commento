package main

import (
	"database/sql"
	"log"
	"time"
)

type SqliteDatabase struct {
	*sql.DB
}

func sqliteInit(params map[string]string) (*SqliteDatabase, error) {
	filename, ok := params["file"]
	if !ok {
		return nil, errorList["err.db.conf.sqlite.filename.missing"]
	}

	db, err := sql.Open("sqlite3", filename)
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

func (db *SqliteDatabase) CreateComment(c *Comment) error {
	statement := `
		SELECT depth, parent FROM comments WHERE rowid=?;
	`
	rows, err := db.Query(statement, c.Parent)
	if err != nil {
		return err
	}
	defer rows.Close()

	depth := 0
	for rows.Next() {
		var pParent int
		if err := rows.Scan(&depth, &pParent); err == nil {
			if depth+1 > 5 {
				c.Parent = pParent
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return err
	}

	statement = `
		INSERT INTO comments(url, name, comment, time, depth, parent) VALUES(?, ?, ?, ?, ?, ?);
	`
	_, err = db.Exec(statement, c.URL, c.Name, c.Comment, time.Now(), depth+1, c.Parent)
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
		c := Comment{}
		if err = rows.Scan(&c.ID, &c.URL, &c.Comment, &c.Name, &c.Timestamp, &c.Parent); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, rows.Err()
}
