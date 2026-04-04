package main

import (
	"net/http"
)

func adminSubscriptionUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken    *string `json:"ownerToken"`
		TargetOwnerHex *string `json:"targetOwnerHex"`
		PlanId        *string `json:"planId"`
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

	if err := subscriptionUpdate(*x.TargetOwnerHex, *x.PlanId); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
