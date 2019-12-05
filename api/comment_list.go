package main

import (
	"database/sql"
	"net/http"
)

func commentList(commenterHex string, domain string, path string, includeUnapproved bool) ([]comment, map[string]commenter, error) {
	// path can be empty
	if commenterHex == "" || domain == "" {
		return nil, nil, errorMissingField
	}

	statement := `
		SELECT
			commentHex,
			commenterHex,
			markdown,
			html,
			parentHex,
			score,
			state,
			deleted,
			creationDate
		FROM comments
		WHERE
			comments.domain = $1 AND
			comments.path = $2
	`

	if !includeUnapproved {
		if commenterHex == "anonymous" {
			statement += `AND state = 'approved'`
		} else {
			statement += `AND (state = 'approved' OR commenterHex = $3)`
		}
	}

	statement += `;`

	var rows *sql.Rows
	var err error

	if !includeUnapproved && commenterHex != "anonymous" {
		rows, err = db.Query(statement, domain, path, commenterHex)
	} else {
		rows, err = db.Query(statement, domain, path)
	}

	if err != nil {
		logger.Errorf("cannot get comments: %v", err)
		return nil, nil, errorInternal
	}
	defer rows.Close()

	commenters := make(map[string]commenter)
	commenters["anonymous"] = commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []comment{}
	for rows.Next() {
		c := comment{}
		if err = rows.Scan(
			&c.CommentHex,
			&c.CommenterHex,
			&c.Markdown,
			&c.Html,
			&c.ParentHex,
			&c.Score,
			&c.State,
			&c.Deleted,
			&c.CreationDate); err != nil {
			return nil, nil, errorInternal
		}

		if commenterHex != "anonymous" {
			statement = `
				SELECT direction
				FROM votes
				WHERE commentHex=$1 AND commenterHex=$2;
			`
			row := db.QueryRow(statement, c.CommentHex, commenterHex)

			if err = row.Scan(&c.Direction); err != nil {
				// TODO: is the only error here that there is no such entry?
				c.Direction = 0
			}
		}

		if commenterHex != c.CommenterHex {
			c.Markdown = ""
		}

		if !includeUnapproved {
			c.State = ""
		}

		comments = append(comments, c)

		if _, ok := commenters[c.CommenterHex]; !ok {
			commenters[c.CommenterHex], err = commenterGetByHex(c.CommenterHex)
			if err != nil {
				logger.Errorf("cannot retrieve commenter: %v", err)
				return nil, nil, errorInternal
			}
		}
	}

	return comments, commenters, nil
}

func commentListHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"CommenterToken"`
		Domain         *string `json:"domain"`
		Path           *string `json:"path"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	path := *x.Path

	d, err := domainGet(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	p, err := pageGet(domain, path)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	commenterHex := "anonymous"
	isModerator := false
	modList := map[string]bool{}

	if *x.CommenterToken != "anonymous" {
		c, err := commenterGetByCommenterToken(*x.CommenterToken)
		if err != nil {
			if err == errorNoSuchToken {
				commenterHex = "anonymous"
			} else {
				bodyMarshal(w, response{"success": false, "message": err.Error()})
				return
			}
		} else {
			commenterHex = c.CommenterHex
		}

		for _, mod := range d.Moderators {
			modList[mod.Email] = true
			if mod.Email == c.Email {
				isModerator = true
			}
		}
	} else {
		for _, mod := range d.Moderators {
			modList[mod.Email] = true
		}
	}

	domainViewRecord(domain, commenterHex)

	comments, commenters, err := commentList(commenterHex, domain, path, isModerator)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	_commenters := map[string]commenter{}
	for commenterHex, cr := range commenters {
		if _, ok := modList[cr.Email]; ok {
			cr.IsModerator = true
		}
		cr.Email = ""
		_commenters[commenterHex] = cr
	}

	bodyMarshal(w, response{
		"success":               true,
		"domain":                domain,
		"comments":              comments,
		"commenters":            _commenters,
		"requireModeration":     d.RequireModeration,
		"requireIdentification": d.RequireIdentification,
		"isFrozen":              d.State == "frozen",
		"isModerator":           isModerator,
		"defaultSortPolicy":     d.DefaultSortPolicy,
		"attributes":            p,
		"configuredOauths": map[string]bool{
			"commento": d.CommentoProvider,
			"google":   googleConfigured && d.GoogleProvider,
			"twitter":  twitterConfigured && d.TwitterProvider,
			"github":   githubConfigured && d.GithubProvider,
			"gitlab":   gitlabConfigured && d.GitlabProvider,
			"sso":      d.SsoProvider,
		},
	})
}
