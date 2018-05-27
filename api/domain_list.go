package main

import (
	"net/http"
)

func domainList(ownerHex string) ([]domain, error) {
	if ownerHex == "" {
		return []domain{}, errorMissingField
	}

	statement := `
    SELECT domain, ownerHex, name, creationDate, state, importedComments, autoSpamFilter, requireModeration, requireIdentification
		FROM domains
		WHERE ownerHex=$1;
	`
	rows, err := db.Query(statement, ownerHex)
	if err != nil {
		logger.Errorf("cannot query domains: %v", err)
		return nil, errorInternal
	}
	defer rows.Close()

	domains := []domain{}
	for rows.Next() {
		d := domain{}
		if err = rows.Scan(&d.Domain, &d.OwnerHex, &d.Name, &d.CreationDate, &d.State, &d.ImportedComments, &d.AutoSpamFilter, &d.RequireModeration, &d.RequireIdentification); err != nil {
			logger.Errorf("cannot Scan domain: %v", err)
			return nil, errorInternal
		}

		d.Moderators, err = domainModeratorList(d.Domain)
		if err != nil {
			return []domain{}, err
		}

		domains = append(domains, d)
	}

	return domains, rows.Err()
}

func domainListHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Session *string `json:"session"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetBySession(*x.Session)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	domains, err := domainList(o.OwnerHex)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true, "domains": domains})
}
