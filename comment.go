package main

import (
	"errors"
	"log"
	"strings"
	"time"
)

// Errors returned when working with comments
var (
	ErrInvalidComment = errors.New("invalid comment, missing required fields")
)

// Comment defines the structure of a commento comment
type Comment struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Comment   string    `json:"comment"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Parent    int       `json:"parent"`
}

func createComment(url string, name string, comment string, parent int) error {
	required := []string{url, name, comment}
	for _, val := range required {
		val = strings.TrimSpace(val)
		if val == "" {
			return ErrInvalidComment
		}
	}

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

	var comments []Comment
	for rows.Next() {
		c := Comment{}
		err := rows.Scan(&c.ID, &c.URL, &c.Comment, &c.Name, &c.Timestamp, &c.Parent)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return comments, nil
}
