package main

import (
	"net/http"
)

func ownerUpdate(ownerHex string, name string) error {
	if ownerHex == "" || name == "" {
		return errorMissingField
	}

	statement := `
		UPDATE owners
		SET name = $2
		WHERE ownerHex = $1;
	`
	_, err := db.Exec(statement, ownerHex, name)
	if err != nil {
		logger.Errorf("cannot update owner: %v", err)
		return errorInternal
	}

	return nil
}

func ownerUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Name       *string `json:"name"`
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

	if err := ownerUpdate(o.OwnerHex, *x.Name); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
