package main

import ()

func pageNew(domain string, path string) error {
	// path can be empty
	if domain == "" {
		return errorMissingField
	}

	statement := `
		INSERT INTO
		pages  (domain, path)
		VALUES ($1,     $2  )
		ON CONFLICT DO NOTHING;
	`
	_, err := db.Exec(statement, domain, path)
	if err != nil {
		logger.Errorf("error inserting new page: %v", err)
		return errorInternal
	}

	return nil
}
