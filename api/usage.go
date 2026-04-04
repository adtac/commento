package main

import (
	"fmt"
	"net/http"
	"time"
)

type usageRecord struct {
	OwnerHex      string `json:"ownerHex"`
	Month         string `json:"month"`
	CommentCount  int    `json:"commentCount"`
	PageViewCount int    `json:"pageViewCount"`
}

func currentMonth() string {
	return time.Now().UTC().Format("2006-01")
}

func usageGet(ownerHex string) (usageRecord, error) {
	if ownerHex == "" {
		return usageRecord{}, errorMissingField
	}

	month := currentMonth()

	statement := `
		SELECT ownerHex, month, commentCount, pageViewCount
		FROM usage
		WHERE ownerHex = $1 AND month = $2;
	`
	row := db.QueryRow(statement, ownerHex, month)

	var u usageRecord
	if err := row.Scan(&u.OwnerHex, &u.Month, &u.CommentCount, &u.PageViewCount); err != nil {
		// No usage row yet means zero usage
		return usageRecord{OwnerHex: ownerHex, Month: month, CommentCount: 0, PageViewCount: 0}, nil
	}

	return u, nil
}

func usageIncrement(ownerHex string, metric string) error {
	if ownerHex == "" {
		return errorMissingField
	}

	month := currentMonth()

	var statement string
	switch metric {
	case "commentCount":
		statement = `
			INSERT INTO usage (ownerHex, month, commentCount, pageViewCount)
			VALUES ($1, $2, 1, 0)
			ON CONFLICT (ownerHex, month)
			DO UPDATE SET commentCount = usage.commentCount + 1;
		`
	case "pageViewCount":
		statement = `
			INSERT INTO usage (ownerHex, month, commentCount, pageViewCount)
			VALUES ($1, $2, 0, 1)
			ON CONFLICT (ownerHex, month)
			DO UPDATE SET pageViewCount = usage.pageViewCount + 1;
		`
	default:
		return fmt.Errorf("unknown usage metric: %s", metric)
	}

	_, err := db.Exec(statement, ownerHex, month)
	if err != nil {
		logger.Errorf("cannot increment usage %s for %s: %v", metric, ownerHex, err)
		return errorInternal
	}

	return nil
}

func usageCheckLimit(ownerHex string, metric string) (bool, error) {
	p, err := planGetForOwner(ownerHex)
	if err != nil {
		return false, err
	}

	u, err := usageGet(ownerHex)
	if err != nil {
		return false, err
	}

	switch metric {
	case "commentCount":
		return u.CommentCount < p.MaxCommentsPerMonth, nil
	case "pageViewCount":
		return u.PageViewCount < p.MaxPageViewsPerMonth, nil
	default:
		return false, fmt.Errorf("unknown usage metric: %s", metric)
	}
}

func domainCountForOwner(ownerHex string) (int, error) {
	statement := `
		SELECT COUNT(*)
		FROM domains
		WHERE ownerHex = $1;
	`
	row := db.QueryRow(statement, ownerHex)

	var count int
	if err := row.Scan(&count); err != nil {
		logger.Errorf("cannot count domains for owner: %v", err)
		return 0, errorInternal
	}

	return count, nil
}

func usageGetHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	u, err := usageGet(o.OwnerHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	p, err := planGetForOwner(o.OwnerHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domainCount, err := domainCountForOwner(o.OwnerHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{
		"success":      true,
		"usage":        u,
		"plan":         p,
		"domainCount":  domainCount,
	})
}
