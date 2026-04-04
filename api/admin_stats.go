package main

import (
	"net/http"
)

func adminStatsHandler(w http.ResponseWriter, r *http.Request) {
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

	var totalOwners int
	row := db.QueryRow(`SELECT COUNT(*) FROM owners;`)
	if err := row.Scan(&totalOwners); err != nil {
		logger.Errorf("cannot count owners: %v", err)
		bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
		return
	}

	var totalDomains int
	row = db.QueryRow(`SELECT COUNT(*) FROM domains;`)
	if err := row.Scan(&totalDomains); err != nil {
		logger.Errorf("cannot count domains: %v", err)
		bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
		return
	}

	var totalComments int
	row = db.QueryRow(`SELECT COUNT(*) FROM comments;`)
	if err := row.Scan(&totalComments); err != nil {
		logger.Errorf("cannot count comments: %v", err)
		bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
		return
	}

	var totalCommenters int
	row = db.QueryRow(`SELECT COUNT(*) FROM commenters;`)
	if err := row.Scan(&totalCommenters); err != nil {
		logger.Errorf("cannot count commenters: %v", err)
		bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
		return
	}

	// Plan distribution
	type planCount struct {
		PlanId string `json:"planId"`
		Count  int    `json:"count"`
	}
	planRows, err := db.Query(`SELECT planId, COUNT(*) FROM subscriptions GROUP BY planId ORDER BY COUNT(*) DESC;`)
	if err != nil {
		logger.Errorf("cannot query plan distribution: %v", err)
		bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
		return
	}
	defer planRows.Close()

	planDistribution := []planCount{}
	for planRows.Next() {
		var pc planCount
		if err := planRows.Scan(&pc.PlanId, &pc.Count); err != nil {
			logger.Errorf("cannot scan plan count: %v", err)
			bodyMarshal(w, response{"success": false, "message": errorInternal.Error()})
			return
		}
		planDistribution = append(planDistribution, pc)
	}

	bodyMarshal(w, response{
		"success":          true,
		"totalOwners":      totalOwners,
		"totalDomains":     totalDomains,
		"totalComments":    totalComments,
		"totalCommenters":  totalCommenters,
		"planDistribution": planDistribution,
	})
}
