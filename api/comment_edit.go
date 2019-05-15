package main

import (
	"net/http"
)

func commentEdit(commentHex string, markdown string) (string, error) {
	if commentHex == "" {
		return "", errorMissingField
	}

	html := markdownToHtml(markdown)

	statement := `
		UPDATE comments
		SET markdown = $2, html = $3
		WHERE commentHex=$1;
	`
	_, err := db.Exec(statement, commentHex, markdown, html)

	if err != nil {
		// TODO: make sure this is the error is actually non-existant commentHex
		return "", errorNoSuchComment
	}

	return html, nil
}

func commentEditHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		CommentHex     *string `json:"commentHex"`
		Markdown       *string `json:"markdown"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	cm, err := commentGetByCommentHex(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if cm.CommenterHex != c.CommenterHex {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	html, err := commentEdit(*x.CommentHex, *x.Markdown)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "html": html})
}
