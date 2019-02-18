package main

import (
	"database/sql"
)

func pageGet(domain string, path string) (page, error) {
	// path can be empty
	if domain == "" {
		return page{}, errorMissingField
	}

	statement := `
		SELECT isLocked, commentCount, stickyCommentHex, title
		FROM pages
		WHERE domain=$1 AND path=$2;
	`
	row := db.QueryRow(statement, domain, path)

	p := page{Domain: domain, Path: path}
	if err := row.Scan(&p.IsLocked, &p.CommentCount, &p.StickyCommentHex, &p.Title); err != nil {
		if err == sql.ErrNoRows {
			// If there haven't been any comments, there won't be a record for this
			// page. The sane thing to do is return defaults.
			// TODO: the defaults are hard-coded in two places: here and the schema
			p.IsLocked = false
			p.CommentCount = 0
			p.StickyCommentHex = "none"
			p.Title = ""
		} else {
			logger.Errorf("error scanning page: %v", err)
			return page{}, errorInternal
		}
	}

	return p, nil
}
