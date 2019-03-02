package main

import (
	"github.com/lib/pq"
	"net/http"
)

func commentCount(domain string, paths []string) (map[string]int, error) {
	commentCounts := map[string]int{}

	if domain == "" {
		return nil, errorMissingField
	}

	if len(paths) == 0 {
		return nil, errorEmptyPaths
	}

	statement := `
		SELECT path, commentCount
		FROM pages
		WHERE domain = $1 AND path = ANY($2);
	`
	rows, err := db.Query(statement, domain, pq.Array(paths))
	if err != nil {
		logger.Errorf("cannot get comments: %v", err)
		return nil, errorInternal
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		var commentCount int
		if err = rows.Scan(&path, &commentCount); err != nil {
			logger.Errorf("cannot scan path and commentCount: %v", err)
			return nil, errorInternal
		}

		commentCounts[path] = commentCount
	}

	return commentCounts, nil
}

func commentCountHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Domain *string   `json:"domain"`
		Paths  *[]string `json:"paths"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)

	commentCounts, err := commentCount(domain, *x.Paths)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commentCounts": commentCounts})
}
