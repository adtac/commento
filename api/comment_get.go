package main

import ()

func commentGetByCommentHex(commentHex string) (comment, error) {
	if commentHex == "" {
		return comment{}, errorMissingField
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
			creationDate
		FROM comments
		WHERE comments.commentHex = $1;
	`
	row := db.QueryRow(statement, commentHex)

	c := comment{}
	if err := row.Scan(
		&c.CommentHex,
		&c.CommenterHex,
		&c.Markdown,
		&c.Html,
		&c.ParentHex,
		&c.Score,
		&c.State,
		&c.CreationDate); err != nil {
		// TODO: is this the only error?
		return c, errorNoSuchComment
	}

	return c, nil
}
