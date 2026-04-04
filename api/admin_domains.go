package main

import (
	"net/http"
)

func adminDomainsListHandler(w http.ResponseWriter, r *http.Request) {
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

	if !o.IsSuperAdmin {
		bodyMarshal(w, response{"success": false, "message": errorNotSuperAdmin.Error()})
		return
	}

	statement := `
		SELECT ` + domainsRowColumns + `
		FROM domains
		ORDER BY creationDate DESC;
	`
	rows, err := db.Query(statement)
	if err != nil {
		logger.Errorf("cannot query all domains: %v", err)
		bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
		return
	}
	defer rows.Close()

	domains := []domain{}
	for rows.Next() {
		var d domain
		if err := domainsRowScan(rows, &d); err != nil {
			logger.Errorf("cannot scan domain: %v", err)
			bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
			return
		}
		domains = append(domains, d)
	}

	bodyMarshal(w, response{"success": true, "domains": domains})
}
