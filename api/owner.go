package main

import (
	"time"
)

type owner struct {
	OwnerHex       string    `json:"ownerHex"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	ConfirmedEmail bool      `json:"confirmedEmail"`
	JoinDate       time.Time `json:"joinDate"`
}
