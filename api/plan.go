package main

import (
	"net/http"
)

type plan struct {
	PlanId               string `json:"planId"`
	Name                 string `json:"name"`
	MaxDomains           int    `json:"maxDomains"`
	MaxCommentsPerMonth  int    `json:"maxCommentsPerMonth"`
	MaxPageViewsPerMonth int    `json:"maxPageViewsPerMonth"`
	IsActive             bool   `json:"isActive"`
}

func planGet(planId string) (plan, error) {
	if planId == "" {
		return plan{}, errorMissingField
	}

	statement := `
		SELECT planId, name, maxDomains, maxCommentsPerMonth, maxPageViewsPerMonth, isActive
		FROM plans
		WHERE planId = $1;
	`
	row := db.QueryRow(statement, planId)

	var p plan
	if err := row.Scan(&p.PlanId, &p.Name, &p.MaxDomains, &p.MaxCommentsPerMonth, &p.MaxPageViewsPerMonth, &p.IsActive); err != nil {
		return plan{}, errorNoSuchPlan
	}

	return p, nil
}

func planList() ([]plan, error) {
	statement := `
		SELECT planId, name, maxDomains, maxCommentsPerMonth, maxPageViewsPerMonth, isActive
		FROM plans
		WHERE isActive = true
		ORDER BY maxDomains ASC;
	`
	rows, err := db.Query(statement)
	if err != nil {
		logger.Errorf("cannot query plans: %v", err)
		return nil, errorInternal
	}
	defer rows.Close()

	plans := []plan{}
	for rows.Next() {
		var p plan
		if err := rows.Scan(&p.PlanId, &p.Name, &p.MaxDomains, &p.MaxCommentsPerMonth, &p.MaxPageViewsPerMonth, &p.IsActive); err != nil {
			logger.Errorf("cannot scan plan: %v", err)
			return nil, errorInternal
		}
		plans = append(plans, p)
	}

	return plans, rows.Err()
}

func planListHandler(w http.ResponseWriter, r *http.Request) {
	plans, err := planList()
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "plans": plans})
}

func planGetForOwner(ownerHex string) (plan, error) {
	sub, err := subscriptionGet(ownerHex)
	if err != nil {
		return plan{}, err
	}

	return planGet(sub.PlanId)
}
