package main

import (
	"time"
)

type commenter struct {
	CommenterHex string    `json:"commenterHex,omitempty"`
	Email        string    `json:"email,omitempty"`
	Name         string    `json:"name"`
	Link         string    `json:"link"`
	Photo        string    `json:"photo"`
	Provider     string    `json:"provider,omitempty"`
	JoinDate     time.Time `json:"joinDate,omitempty"`
	IsModerator  bool      `json:"isModerator"`
}
