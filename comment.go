package main

import (
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Comment   string    `json:"comment"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Parent    int       `json:"parent"`
}

func createComment(url string, name string, comment string, parent int) error {
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
				parent = pParent;
			}
		}
	}

	statement = `
		INSERT INTO comments(url, name, comment, time, depth, parent) VALUES(?, ?, ?, ?, ?, ?);
	`
	_, err = db.Exec(statement, url, name, comment, time.Now(), depth+1, parent)
	return err
}

func getComments(url string) ([]Comment, error) {
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

	return comments, nil
}
