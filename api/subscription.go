package main

import (
	"net/http"
	"time"
)

type subscription struct {
	OwnerHex           string    `json:"ownerHex"`
	PlanId             string    `json:"planId"`
	Status             string    `json:"status"`
	CurrentPeriodStart time.Time `json:"currentPeriodStart"`
	CurrentPeriodEnd   time.Time `json:"currentPeriodEnd"`
	CreatedDate        time.Time `json:"createdDate"`
}

func subscriptionGet(ownerHex string) (subscription, error) {
	if ownerHex == "" {
		return subscription{}, errorMissingField
	}

	statement := `
		SELECT ownerHex, planId, status, currentPeriodStart, currentPeriodEnd, createdDate
		FROM subscriptions
		WHERE ownerHex = $1;
	`
	row := db.QueryRow(statement, ownerHex)

	var s subscription
	if err := row.Scan(&s.OwnerHex, &s.PlanId, &s.Status, &s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CreatedDate); err != nil {
		return subscription{}, errorNoSuchSubscription
	}

	return s, nil
}

func subscriptionNew(ownerHex string, planId string) error {
	if ownerHex == "" || planId == "" {
		return errorMissingField
	}

	now := time.Now().UTC()
	statement := `
		INSERT INTO subscriptions (ownerHex, planId, status, currentPeriodStart, currentPeriodEnd, createdDate)
		VALUES ($1, $2, 'active', $3, $4, $3)
		ON CONFLICT (ownerHex) DO UPDATE SET planId = $2, status = 'active', currentPeriodStart = $3, currentPeriodEnd = $4;
	`
	_, err := db.Exec(statement, ownerHex, planId, now, now.AddDate(100, 0, 0))
	if err != nil {
		logger.Errorf("cannot insert subscription: %v", err)
		return errorInternal
	}

	return nil
}

func subscriptionUpdate(ownerHex string, planId string) error {
	if ownerHex == "" || planId == "" {
		return errorMissingField
	}

	// Verify plan exists
	if _, err := planGet(planId); err != nil {
		return err
	}

	statement := `
		UPDATE subscriptions
		SET planId = $2
		WHERE ownerHex = $1;
	`
	_, err := db.Exec(statement, ownerHex, planId)
	if err != nil {
		logger.Errorf("cannot update subscription: %v", err)
		return errorInternal
	}

	return nil
}

func subscriptionGetHandler(w http.ResponseWriter, r *http.Request) {
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

	sub, err := subscriptionGet(o.OwnerHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	p, err := planGet(sub.PlanId)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "subscription": sub, "plan": p})
}

func subscriptionUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		PlanId     *string `json:"planId"`
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

	if err := subscriptionUpdate(o.OwnerHex, *x.PlanId); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
