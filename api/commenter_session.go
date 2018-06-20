package main

import (
	"time"
)

// A session is a 3-field entry of a token, a hex, and a creation date. Do
// not confuse session and token; the token is just an identifying string,
// while the session contains more information.
type commenterSession struct {
	CommenterToken string    `json:"commenterToken"`
	CommenterHex   string    `json:"commenterHex"`
	CreationDate   time.Time `json:"creationDate"`
}
