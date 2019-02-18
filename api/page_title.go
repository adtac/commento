package main

import ()

func pageTitleUpdate(domain string, path string) (string, error) {
	title, err := htmlTitleGet("http://" + domain + path)
	if err != nil {
		// This could fail due to a variety of reasons that we can't control such
		// as the user's URL 404 or something, so let's not pollute the error log
		// with messages. Just use a sane title. Maybe we'll have the ability to
		// retry later.
		logger.Errorf("%v", err)
		title = domain
	}

	statement := `
		UPDATE pages
		SET title = $3
		WHERE domain = $1 AND path = $2;
	`
	_, err = db.Exec(statement, domain, path, title)
	if err != nil {
		logger.Errorf("cannot update pages table with title: %v", err)
		return "", err
	}

	return title, nil
}
