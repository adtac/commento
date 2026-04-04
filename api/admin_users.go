package main

import (
	"net/http"
)

func adminUsersListHandler(w http.ResponseWriter, r *http.Request) {
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
		SELECT ` + ownersRowColumns + `
		FROM owners
		ORDER BY joinDate DESC;
	`
	rows, err := db.Query(statement)
	if err != nil {
		logger.Errorf("cannot query owners: %v", err)
		bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
		return
	}
	defer rows.Close()

	owners := []owner{}
	for rows.Next() {
		var o owner
		if err := ownersRowScan(rows, &o); err != nil {
			logger.Errorf("cannot scan owner: %v", err)
			bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
			return
		}
		owners = append(owners, o)
	}

	bodyMarshal(w, response{"success": true, "owners": owners})
}

func adminUsersGetHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		OwnerHex   *string `json:"ownerHex"`
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

	target, err := ownerGetByOwnerHex(*x.OwnerHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	sub, err := subscriptionGet(*x.OwnerHex)
	if err != nil {
		sub = subscription{PlanId: "free"}
	}

	usage, err := usageGet(*x.OwnerHex)
	if err != nil {
		usage = usageRecord{}
	}

	domainCount, err := domainCountForOwner(*x.OwnerHex)
	if err != nil {
		domainCount = 0
	}

	bodyMarshal(w, response{
		"success":      true,
		"owner":        target,
		"subscription": sub,
		"usage":        usage,
		"domainCount":  domainCount,
	})
}
